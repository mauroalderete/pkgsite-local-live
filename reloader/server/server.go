package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mauroalderete/pkgsite-local-live/reloader/interceptor/livereload"
	"github.com/mauroalderete/pkgsite-local-live/reloader/reverseproxy"
)

type server struct {
	origin            *url.URL
	public            *url.URL
	reloadSnippetPath string
	proxy             *reverseproxy.ReverseProxy
}

func (s *server) Run() error {

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("must to upgrade connection")
	})

	serverMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("must to redirect")
		s.proxy.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(s.public.Host, serverMux)
	if err != nil {
		return fmt.Errorf("failed to execute the main server: %v", err)
	}

	return nil
}

type Configurator interface {
	Origin(address string) error
	Public(address string) error
	ReloadSnippet(path string) error
}

type configure struct {
	pool []func(*server) error
}

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

	return srv, nil
}
