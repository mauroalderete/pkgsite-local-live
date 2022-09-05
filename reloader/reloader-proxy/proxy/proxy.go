// Package proxy allows instance a simple reverse proxy with capacity to implement many interceptors.
package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	neturl "net/url"

	"github.com/mauroalderete/pkgsite-local-live/reloader/reloader-proxy/interceptor"
)

// proxy execute a reverse proxy and manage the interceptors configured.
type proxy struct {

	// origin is the backend endpoint that the proxy query by each request of the client.
	origin *neturl.URL

	// endpoint is the frontend endpoint for clients to access.
	endpoint *neturl.URL

	// upgrade is the backend endpoint to connections of upgrade type incoming.
	upgrade *neturl.URL

	// proxy is the httputil.ReverseProxy instance that is executed.
	proxy *httputil.ReverseProxy

	// interceptors is a list of the all interceptor.Interceptor configured.
	interceptors map[string]interceptor.Interceptor
}

func (rp *proxy) director(request *http.Request) {
	log.Printf("[-------------------------------------------]\n")
	// log.Printf("\treq %v\n", req.ContentLength)
	// log.Printf("\treq %v\n", req.Header)
	// log.Printf("\treq %v\n", req.Host)
	// log.Printf("\treq %v\n", req.Method)
	// log.Printf("\treq %v\n", req.Proto)
	// log.Printf("\treq %v\n", req.RemoteAddr)
	// log.Printf("\treq %v\n", req.RequestURI)
	// log.Printf("\treq %v\n", req.URL)
	// log.Printf("\treq %v\n", req.URL.Host)
	// log.Printf("\treq %v\n", req.URL.Path)
	// log.Printf("\treq %v\n", req.URL.Scheme)
	// log.Printf("\treq %v\n", req.URL.User)
	// log.Printf("\treq %v\n", req.URL.Opaque)

	log.Printf("[director] before %v %v %v\n", request.Host, request.URL.Host, request.URL.Scheme)

	if request.Header.Get("Connection") == "Upgrade" &&
		request.Header.Get("Upgrade") == "websocket" &&
		rp.upgrade != nil {
		redirectTo(request, *rp.upgrade)
		log.Printf("[director] after by upgrade %v %v %v\n", request.Host, request.URL.Host, request.URL.Scheme)
		return
	}

	redirectTo(request, *rp.origin)
	log.Printf("[director] after by default %v %v %v\n", request.Host, request.URL.Host, request.URL.Scheme)
}

func redirectTo(request *http.Request, target neturl.URL) {
	request.Host = target.Host
	request.URL.Host = target.Host
	request.URL.Scheme = target.Scheme
	request.RequestURI = ""
}

// modify iterates for each interceptor and executes his handler if needed.
func (rp *proxy) modify(r *http.Response) error {

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

// Run start to lisent and serve the reverse proxy
func (rp *proxy) Run() error {
	err := http.ListenAndServe(rp.endpoint.Host, rp.proxy)
	if err != nil {
		return fmt.Errorf("reloader proxy failed: %v", err)
	}
	return nil
}

// Configurer defines the available options to configure a new instance of proxy.proxy
type Configurer interface {

	// Origin allows set the endpoint backend url
	Origin(url string) error

	// Endpoint allows set the endpoint frontend url
	Endpoint(url string) error

	// Upgrade allows set the endpoint to connections upgrade type.
	Upgrade(url string) error

	// AddInterceptor allows loading a new interceptor that the proxy must be execute by each request.
	//
	// Receives a name to identify the interceptor loaded.
	AddInterceptor(name string, interceptor interceptor.Interceptor) error
}

// configurerPool implements proxy.Configurer
type configurerPool struct {
	pool []func(*proxy) error
}

// Origin implements proxy.Configurer.Origin method
func (c *configurerPool) Origin(url string) error {

	addr, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse origin url: %v", err)
	}

	c.pool = append(c.pool, func(rp *proxy) error {
		rp.origin = addr
		return nil
	})

	return nil
}

// Endpoint implements proxy.Configurer.Endpoint method
func (c *configurerPool) Endpoint(url string) error {

	addr, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint url: %v", err)
	}

	c.pool = append(c.pool, func(rp *proxy) error {
		rp.endpoint = addr
		return nil
	})

	return nil
}

// Upgrade implements proxy.Configurer.Upgrade method
func (c *configurerPool) Upgrade(url string) error {

	addr, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse origin url: %v", err)
	}

	c.pool = append(c.pool, func(rp *proxy) error {
		rp.upgrade = addr
		return nil
	})

	return nil
}

// AddInterceptor implements proxy.Configurer.AddInterceptor method
func (c *configurerPool) AddInterceptor(name string, interceptor interceptor.Interceptor) error {

	c.pool = append(c.pool, func(rp *proxy) error {

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
func New(options ...func(Configurer) error) (*proxy, error) {

	configurer := &configurerPool{}

	for _, option := range options {
		err := option(configurer)
		if err != nil {
			return nil, fmt.Errorf("failed to load options: %v", err)
		}
	}

	proxy := &proxy{
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
