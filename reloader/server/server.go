// Package server implements a basic server HTTP to expose routes
// that allows to handle a reverse proxy and a websocket to manage the livereload feature.
package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/mauroalderete/pkgsite-local-live/reloader/interceptor/livereload"
	"github.com/mauroalderete/pkgsite-local-live/reloader/reverseproxy"
	"github.com/mauroalderete/pkgsite-local-live/reloader/websocketserver"
)

// server allow initialize the HTTP server that is manager to replicate the pkgsite endpoint
// and serve a websocket connection to handle the livereload system.
//
// Stores the address to origin and public endpoints.
// Initialize a instance of reverseproxy.ReverseProxy and websocketserver.WebsocketServer.
type server struct {
	origin            *url.URL
	public            *url.URL
	reloadSnippetPath string
	proxy             *reverseproxy.ReverseProxy
	websocket         *websocketserver.WebsocketServer
}

// Run uploads a new serverMux and launch it.
func (s *server) Run() error {

	serverMux := http.NewServeMux()

	// handler to accept a new websocket connection
	serverMux.HandleFunc("/ws", func(response http.ResponseWriter, request *http.Request) {
		s.websocket.WebsocketHandler(response, request)
	})

	// handler to send broadcast reload signal
	serverMux.HandleFunc("/ws/reload", func(response http.ResponseWriter, request *http.Request) {
		s.websocket.ReloadHandler(response, request)
	})

	// handler to redirect any connection
	serverMux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		s.proxy.ServeHTTP(response, request)
	})

	err := http.ListenAndServe(s.public.Host, serverMux)
	if err != nil {
		return fmt.Errorf("failed to execute the main server: %v", err)
	}

	return nil
}

// Configurator defines the properties configurables to instance a new Server
type Configurator interface {
	// Origin allows set the address to the origin endpoint of the reverse proxy.
	Origin(address string) error

	// Public allows set the address to is expose the endpoint of the reverse proxy.
	Public(address string) error

	// ReloadSnippet allows set the path to the file that contains the snippet
	// that is needed to inject in each request with html content
	// to the browser can be reloaded when it needed.
	ReloadSnippet(path string) error
}

// Implement server.Configurator interface. Stores a pool of configurations callback
type configure struct {
	pool []func(*server) error
}

// Origin implement server.Configurator.Origin method
func (c *configure) Origin(address string) error {
	u, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("failed to parse origin address '%s': %v", address, err)
	}

	c.pool = append(c.pool, func(s *server) error {
		s.origin = u
		return nil
	})

	return nil
}

// Public implement server.Configurator.Public method
func (c *configure) Public(address string) error {

	u, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("failed to parse public address '%s': %v", address, err)
	}

	c.pool = append(c.pool, func(s *server) error {
		s.public = u
		return nil
	})

	return nil
}

// ReloadSnippet implement server.Configurator.ReloadSnippet method
func (c *configure) ReloadSnippet(path string) error {

	if path == "" {
		return fmt.Errorf("reload snippet path cannot be empty")
	}

	c.pool = append(c.pool, func(s *server) error {
		s.reloadSnippetPath = path
		return nil
	})

	return nil
}

// New instances of a new server object using the properties configured through the callbacks options list.
//
// If the options are accepted, loads a new instances of reverseproxy.ReverseProxy,
// a livereload.Livereload interceptor and a websocketserver.WebsockerServer manager.
func New(options ...func(Configurator) error) (*server, error) {

	cnf := &configure{}

	for _, option := range options {
		err := option(cnf)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare option to configure: %v", err)
		}
	}

	srv := &server{}

	for _, config := range cnf.pool {
		err := config(srv)
		if err != nil {
			return nil, fmt.Errorf("failed to apply configuration: %v", err)
		}
	}

	if srv.origin == nil {
		return nil, fmt.Errorf("origin address is required")
	}

	if srv.public == nil {
		return nil, fmt.Errorf("public address is required")
	}

	// load a reverse proxy instance
	rp, err := reverseproxy.New(func(c reverseproxy.Configurer) error {
		err := c.Origin(srv.origin.String())
		if err != nil {
			return fmt.Errorf("failed to set the origin address of the reverse proxy: %v", err)
		}

		err = c.Endpoint(srv.public.String())
		if err != nil {
			return fmt.Errorf("failed to set the public address of the reverse proxy: %v", err)
		}

		livereload, err := livereload.New(func(c livereload.Configurer) error {
			err := c.UpgradeEndpoint(srv.public.String() + "/ws")
			if err != nil {
				return fmt.Errorf("failed to set the upgrade endpoint to livereload interceptor: %v", err)
			}

			err = c.WebserviceInjectable(srv.reloadSnippetPath)
			if err != nil {
				return fmt.Errorf("failed to set the reload snippet path to livereload interceptor: %v", err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to set livereload interceptor of the reverse proxy: %v", err)
		}

		c.AddInterceptor("livereload", livereload)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to up the reverse proxy: %v", err)
	}

	srv.proxy = rp

	// prepare a websocket server
	ws, err := websocketserver.New(func(c websocketserver.Configurator) error {
		err = c.Endpoint(srv.public.String())
		if err != nil {
			return fmt.Errorf("failed to set endpoint to websocket server: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to up the websocket server: %v", err)
	}

	srv.websocket = ws

	return srv, nil
}
