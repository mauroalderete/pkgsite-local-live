package reloaderwebsocket

import (
	"fmt"
	"log"
	"net/http"
	neturl "net/url"

	"github.com/gorilla/websocket"
)

var connectionCount = 0

type reloaderWebsocket struct {
	endpoint *neturl.URL
	server   *http.ServeMux
}

func (rw *reloaderWebsocket) websocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[echo] start\n")

	connectionCount++
	var connectionID = connectionCount
	log.Printf("[echo] connections count: %d\n", connectionCount)

	upgrader := websocket.Upgrader{}

	log.Printf("[echo %d] creating ws\n", connectionID)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[echo %d] failed upgrade connection: %v\n", connectionID, err)
		return
	}
	defer c.Close()

	log.Printf("[echo %d] start processing...\n", connectionID)
	for {
		log.Printf("[echo %d] wait for message\n", connectionID)
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("[echo %d] failed read a message: %v", connectionID, err)
			break
		}
		log.Printf("[echo %d] recived: %s of type %d\n", connectionID, string(message), mt)

		log.Printf("[echo %d] sending\n", connectionID)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Printf("[echo %d] failed to send message: %v", connectionID, err)
			break
		}
	}
	log.Printf("[echo %d] end\n", connectionID)
}

func (rw *reloaderWebsocket) reloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Emit broadcasting signal to reload\n")
}

func (rw *reloaderWebsocket) Run() {
	log.Printf("[reloadwebsocket] Running at %s\n", rw.endpoint.String())
	err := http.ListenAndServe(rw.endpoint.String(), rw.server)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

type ConfigurerNew interface {
	Endpoint(url string) error
}

type configurerPoolNew struct {
	pool []func(*reloaderWebsocket) error
}

func (c *configurerPoolNew) Endpoint(url string) error {

	endpoint, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse origin url: %v", err)
	}

	c.pool = append(c.pool, func(rw *reloaderWebsocket) error {
		rw.endpoint = endpoint
		return nil
	})

	return nil
}

func New(options ...func(ConfigurerNew) error) (*reloaderWebsocket, error) {
	configurer := &configurerPoolNew{}

	for _, option := range options {
		err := option(configurer)
		if err != nil {
			return nil, fmt.Errorf("failed to load options: %v", err)
		}
	}

	websocket := &reloaderWebsocket{}

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

	websocket.server.HandleFunc("/", websocket.websocketHandler)
	websocket.server.HandleFunc("/reload", websocket.reloadHandler)

	return websocket, nil
}
