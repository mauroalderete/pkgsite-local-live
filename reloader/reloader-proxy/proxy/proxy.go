package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	neturl "net/url"

	"github.com/mauroalderete/pkgsite-local-live/reloader/reloader-proxy/interceptor"
)

type proxy struct {
	origin       *neturl.URL
	endpoint     *neturl.URL
	proxy        *httputil.ReverseProxy
	interceptors map[string]interceptor.Interceptor
}

func (rp *proxy) director(req *http.Request) {
	req.Host = rp.origin.Host
	req.URL.Host = rp.origin.Host
	req.URL.Scheme = rp.origin.Scheme
	req.RequestURI = ""
}

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

func (rp *proxy) Run() error {
	err := http.ListenAndServe(rp.endpoint.Host, rp.proxy)
	if err != nil {
		return fmt.Errorf("reloader proxy failed: %v", err)
	}
	return nil
}

type ConfigurerNew interface {
	SetOrigin(url string) error
	SetEndpoint(url string) error
	AddInterceptor(name string, interceptor interceptor.Interceptor) error
}

type configurerPoolNew struct {
	pool []func(*proxy) error
}

func (c *configurerPoolNew) SetOrigin(url string) error {

	o, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse origin url: %v", err)
	}

	c.pool = append(c.pool, func(rp *proxy) error {
		rp.origin = o
		return nil
	})

	return nil
}

func (c *configurerPoolNew) SetEndpoint(url string) error {

	o, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint url: %v", err)
	}

	c.pool = append(c.pool, func(rp *proxy) error {
		rp.endpoint = o
		return nil
	})

	return nil
}

func (c *configurerPoolNew) AddInterceptor(name string, interceptor interceptor.Interceptor) error {

	c.pool = append(c.pool, func(rp *proxy) error {

		if _, ok := rp.interceptors[name]; ok {
			return fmt.Errorf("failed to load an new interceptor: it already exists an interceptor named %s", name)
		}

		rp.interceptors[name] = interceptor

		return nil
	})
	return nil
}

func New(options ...func(ConfigurerNew) error) (*proxy, error) {

	configurer := &configurerPoolNew{}

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
