// Package websocketserver contains a websocket server that support multiple connections to send a reload signal.
package websocketserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"

	"github.com/mauroalderete/pkgsite-local-live/websocketconnections"
)

// WebsocketServer stores a server instance and the connections establishment
type WebsocketServer struct {
	endpoint    *neturl.URL
	server      *http.ServeMux
	connections map[string]*websocketconnections.Connection
}

// responseError writes an error message and print it although standar logger.
//
// This method is used to simplify the error response when a request process failed.
func (rw *WebsocketServer) responseError(w io.Writer, message error) {
	log.Println(message.Error())

	_, err := io.WriteString(w, message.Error())
	if err != nil {
		log.Printf("failed to send a response to requester: %v", err)
	}
}

// WebsocketHandler handle the request to upgrade it and establishment a web socket connection.
//
// Creates a new websocket connection from request and response objects.
// The new connection is stored an internal list to maintain live the comunication
// and have the connection when the reload signal will be broadcasting.
//
// This method is called if the server executes the Run method,
// or can be launched manually if you pass the correct arguments.
func (rw *WebsocketServer) WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	// Creates a new websocket connection
	connection, err := websocketconnections.New(func(c websocketconnections.Configurer) error {
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

	// Opens the websocket connection
	err = connection.Open()
	if err != nil {
		rw.responseError(w, fmt.Errorf("failed to create a connection: %v", err))
		return
	}

	// Stores the websocket connection to send reload signal later
	rw.connections[connection.UUID()] = connection

	// Runs the websocket connection and wait to ends.
	err = connection.Start()
	if err != nil {
		rw.responseError(w, fmt.Errorf("failed to start a connection: %v", err))
		return
	}

	// Removes the connection terminated from the list
	delete(rw.connections, connection.UUID())

	// Closes the connection if it isn't yet.
	defer connection.Close()
}

// ReloadHandler sends reload signal to all connections stored.
//
// The arguments aren't used, but it is maintain to compatibility with [http.Handler] interface.
func (rw *WebsocketServer) ReloadHandler(w http.ResponseWriter, r *http.Request) {
	for _, conn := range rw.connections {
		log.Printf("send reload signal to %s connection\n", conn.UUID())
		conn.Reload()
	}
}

// Run starts to listen and serve the current server on address configured.
//
// This methods is blocked. Returns an error if something was wrong when started or during the execution.
func (rw *WebsocketServer) Run() error {
	log.Printf("Websocket server running at %s\n", rw.endpoint.String())
	err := http.ListenAndServe(rw.endpoint.String(), rw.server)
	if err != nil {
		return fmt.Errorf("websocket server failed: %v", err)
	}
	return nil
}

// Stop allows stop all connections.
func (rw *WebsocketServer) Stop() {
	for _, conn := range rw.connections {
		conn.Stop()
	}
}

// Configurator defines the optionable configurations to instance a new WebsocketServer.
type Configurator interface {

	// Endpoint allows set the endpoint address of the websocket.
	Endpoint(url string) error
}

// configurer implements [websocketserver.Configurator]. Maintains a pool with configurations to execute.
type configurer struct {
	pool []func(*WebsocketServer) error
}

// Endpoint implements [websocketserver.Configurator.Endpoint] method.
func (c *configurer) Endpoint(url string) error {

	endpoint, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint url: %v", err)
	}

	c.pool = append(c.pool, func(rw *WebsocketServer) error {
		rw.endpoint = endpoint
		return nil
	})

	return nil
}

// New returns a new [websocketserver.WebsocketServer] instance with the endpoint set.
//
// Initializes a [http.ServerMux] with the two routes to handle new websocket connections and reload signal.
func New(options ...func(Configurator) error) (*WebsocketServer, error) {
	configurer := &configurer{}

	for _, option := range options {
		err := option(configurer)
		if err != nil {
			return nil, fmt.Errorf("failed to load options: %v", err)
		}
	}

	websocket := &WebsocketServer{}

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

	websocket.connections = make(map[string]*websocketconnections.Connection)
	websocket.server.HandleFunc("/", websocket.WebsocketHandler)
	websocket.server.HandleFunc("/reload", websocket.ReloadHandler)

	return websocket, nil
}
