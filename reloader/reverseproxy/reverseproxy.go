// Package reverseproxy allows instance a simple reverse proxy with capacity to implement many interceptors.
package reverseproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/mauroalderete/pkgsite-local-live/reloader/interceptor"
)

// ReverseProxy execute a reverse ReverseProxy and manage the interceptors configured.
type ReverseProxy struct {

	// origin is the backend endpoint that the proxy query by each request of the client.
	origin *url.URL

	// endpoint is the frontend endpoint for clients to access.
	endpoint *url.URL

	// proxy is the httputil.ReverseProxy instance that is executed.
	proxy *httputil.ReverseProxy

	// interceptors is a list of the all interceptor.Interceptor configured.
	interceptors map[string]interceptor.Interceptor
}

func (rp *ReverseProxy) director(request *http.Request) {
	redirectTo(request, *rp.origin)
}

func redirectTo(request *http.Request, target url.URL) {
	request.Host = target.Host
	request.URL.Host = target.Host
	request.URL.Scheme = target.Scheme
	request.RequestURI = ""
}

// modify iterates for each interceptor and executes his handler if needed.
func (rp *ReverseProxy) modify(r *http.Response) error {

	// iterates by each interceptor configured to check if the rules are passed.
	// In this case, executes the correspondent interceptor.
	for name, interceptor := range rp.interceptors {
		accepted := true
		for _, rule := range interceptor.Rules() {
			if !rule(r) {
				accepted = false
				break
			}
		}

		if !accepted {
			break
		}

		handler := interceptor.Handler()
		err := handler(r)
		if err != nil {
			return fmt.Errorf("interceptor '%s' failed to run: %v", name, err)
		}
	}

	return nil
}

// Run starts to lisent and serve the reverse proxy
func (rp *ReverseProxy) Run() error {
	err := http.ListenAndServe(rp.endpoint.Host, rp.proxy)
	if err != nil {
		return fmt.Errorf("reloader proxy failed: %v", err)
	}
	return nil
}

// ServeHTTP allows execute a request parse manually
//
// Receives the request data that the reverse proxy handle to apply the correspondent redirection.
func (rp *ReverseProxy) ServeHTTP(response http.ResponseWriter, request *http.Request) error {
	rp.proxy.ServeHTTP(response, request)
	return nil
}

// Configurer defines the available options to configure a new instance of proxy.proxy
type Configurer interface {

	// Origin allows set the endpoint backend url
	Origin(address string) error

	// Public allows set the endpoint frontend url
	Public(address string) error

	// AddInterceptor allows loading a new interceptor that the proxy must be execute by each request.
	//
	// Receives a name to identify the interceptor loaded.
	AddInterceptor(name string, interceptor interceptor.Interceptor) error
}

// configurerPool implements proxy.Configurer
type configurerPool struct {
	pool []func(*ReverseProxy) error
}

// Origin implements proxy.Configurer.Origin method
func (c *configurerPool) Origin(address string) error {

	addr, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("failed to parse origin url: %v", err)
	}

	c.pool = append(c.pool, func(rp *ReverseProxy) error {
		rp.origin = addr
		return nil
	})

	return nil
}

// Endpoint implements proxy.Configurer.Endpoint method
func (c *configurerPool) Public(address string) error {

	addr, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint url: %v", err)
	}

	c.pool = append(c.pool, func(rp *ReverseProxy) error {
		rp.endpoint = addr
		return nil
	})

	return nil
}

// AddInterceptor implements proxy.Configurer.AddInterceptor method
func (c *configurerPool) AddInterceptor(name string, interceptor interceptor.Interceptor) error {

	c.pool = append(c.pool, func(rp *ReverseProxy) error {

		if _, ok := rp.interceptors[name]; ok {
			return fmt.Errorf("failed to load an new interceptor: it already exists an interceptor named %s", name)
		}

		rp.interceptors[name] = interceptor

		return nil
	})
	return nil
}

// New returns a new proxy instaceconfigured
//
// Receives a list of options callback with the configurations to apply.
func New(options ...func(Configurer) error) (*ReverseProxy, error) {

	configurer := &configurerPool{}

	for _, option := range options {
		err := option(configurer)
		if err != nil {
			return nil, fmt.Errorf("failed to load options: %v", err)
		}
	}

	proxy := &ReverseProxy{
		interceptors: make(map[string]interceptor.Interceptor),
	}

	proxy.proxy = &httputil.ReverseProxy{
		Director:       proxy.director,
		ModifyResponse: proxy.modify,
		ErrorLog:       log.Default(),
	}

	for _, config := range configurer.pool {
		err := config(proxy)
		if err != nil {
			return nil, fmt.Errorf("failed to apply options: %v", err)
		}
	}

	if proxy.origin == nil {
		return nil, fmt.Errorf("origin is required")
	}

	if proxy.endpoint == nil {
		return nil, fmt.Errorf("endpoint is required")
	}

	return proxy, nil
}
