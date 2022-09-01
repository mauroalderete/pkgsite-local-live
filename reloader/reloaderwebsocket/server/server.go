package server

import (
	"fmt"
	"log"
	"net/http"
	neturl "net/url"
	"strings"

	"github.com/gorilla/websocket"
)

var connectionCount = 0

type server struct {
	endpoint     *neturl.URL
	server       *http.ServeMux
	reloadSignal chan bool
	stopSignal   chan bool
}

func (rw *server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[reloadwebsocket] start\n")

	connectionCount++
	var connectionID = connectionCount
	log.Printf("[reloadwebsocket] connections count: %d\n", connectionCount)

	upgrader := websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		origin := r.Header.Values("Origin")
		if len(origin) != 1 {
			return false
		}
		return strings.HasPrefix(origin[0], "http://localhost")
	}

	log.Printf("[reloadwebsocket %d] creating ws\n", connectionID)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[reloadwebsocket %d] failed upgrade connection: %v\n", connectionID, err)
		return
	}
	defer log.Printf("reloadwebsocket %d] defering...", connectionID)
	defer c.Close()

	stopReloadHandler := make(chan bool)

	c.SetCloseHandler(func(code int, text string) error {
		log.Printf("[reloadwebsocket %d] clossing... %d %s", connectionID, code, text)
		return nil
	})

	go func() {
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("[reloadwebsocket %d] failed read %v\n", connectionID, err)
				//<-
				return
			}
			log.Printf("[reloadwebsocket %d] recibe %v %v\n", connectionID, mt, message)
		}
	}()

	go func() {
		log.Printf("[reloadwebsocket %d] start processing...\n", connectionID)

		for {
			select {
			case <-rw.reloadSignal:
				{
					log.Printf("[reloadwebsocket %d] reload signal", connectionID)
					err := c.WriteMessage(1, []byte("reload"))
					if err != nil {
						log.Printf("[reloadwebsocket %d] failed to send message: %v", connectionID, err)
						break
					}
				}
			case <-rw.stopSignal:
				{
					log.Printf("[reloadwebsocket %d] stop signal", connectionID)
					stopReloadHandler <- true
				}
			}
		}
	}()

	<-stopReloadHandler
}

func (rw *server) reloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Emit broadcasting signal to reload\n")
	rw.reloadSignal <- true
}

func (rw *server) Run() {
	log.Printf("[reloadwebsocket] Running at %s\n", rw.endpoint.String())
	err := http.ListenAndServe(rw.endpoint.String(), rw.server)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func (rw *server) Stop() {
	rw.stopSignal <- true
}

type ConfigurerNew interface {
	Endpoint(url string) error
}

type configurerPoolNew struct {
	pool []func(*server) error
}

func (c *configurerPoolNew) Endpoint(url string) error {

	endpoint, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse origin url: %v", err)
	}

	c.pool = append(c.pool, func(rw *server) error {
		rw.endpoint = endpoint
		return nil
	})

	return nil
}

func New(options ...func(ConfigurerNew) error) (*server, error) {
	configurer := &configurerPoolNew{}

	for _, option := range options {
		err := option(configurer)
		if err != nil {
			return nil, fmt.Errorf("failed to load options: %v", err)
		}
	}

	websocket := &server{}

	for _, config := range configurer.pool {
		err := config(websocket)
		if err != nil {
			return nil, fmt.Errorf("failed to apply options: %v", err)
		}
	}

	if websocket.endpoint == nil {
		return nil, fmt.Errorf("endpoint is required")
	}

	websocket.server = http.NewServeMux()

	websocket.reloadSignal = make(chan bool)
	websocket.stopSignal = make(chan bool)

	websocket.server.HandleFunc("/", websocket.websocketHandler)
	websocket.server.HandleFunc("/reload", websocket.reloadHandler)

	return websocket, nil
}
