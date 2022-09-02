package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"

	"github.com/mauroalderete/pkgsite-local-live/reloader/reloaderwebsocket/connections"
)

type server struct {
	endpoint    *neturl.URL
	server      *http.ServeMux
	connections map[string]*connections.Connection
}

func (rw *server) responseError(w io.Writer, message error) {
	log.Println(message.Error())

	_, err := io.WriteString(w, message.Error())
	if err != nil {
		log.Printf("failed to send a response to requester: %v", err)
	}
}

func (rw *server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[reloadwebsocket] start\n")

	connection, err := connections.New(func(c connections.Configurer) error {
		err := c.Request(r)
		if err != nil {
			return fmt.Errorf("failed to config request: %v", err)
		}

		err = c.ResponseWriter(w)
		if err != nil {
			return fmt.Errorf("failed to config response: %v", err)
		}

		return nil
	})
	if err != nil {
		rw.responseError(w, fmt.Errorf("failed to create a connection: %v", err))
		return
	}

	err = connection.Open()
	if err != nil {
		rw.responseError(w, fmt.Errorf("failed to create a connection: %v", err))
		return
	}

	rw.connections[connection.UUID()] = connection

	err = connection.Start()
	if err != nil {
		rw.responseError(w, fmt.Errorf("failed to start a connection: %v", err))
		return
	}

	delete(rw.connections, connection.UUID())

	defer connection.Close()
}

func (rw *server) reloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Emit broadcasting signal to reload\n")

	for _, conn := range rw.connections {
		log.Printf("send reload signal to %s connection\n", conn.UUID())
		conn.Reload()
	}
}

func (rw *server) Run() {
	log.Printf("[reloadwebsocket] Running at %s\n", rw.endpoint.String())
	err := http.ListenAndServe(rw.endpoint.String(), rw.server)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func (rw *server) Stop() {
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

	websocket.connections = make(map[string]*connections.Connection)
	websocket.server.HandleFunc("/", websocket.websocketHandler)
	websocket.server.HandleFunc("/reload", websocket.reloadHandler)

	return websocket, nil
}
