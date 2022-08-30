package reloaderproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	neturl "net/url"
)

type reloaderProxy struct {
	origin       *neturl.URL
	endpoint     *neturl.URL
	proxy        *httputil.ReverseProxy
	interceptors map[string]func(*http.Response) error
}

func (rp *reloaderProxy) director(req *http.Request) {
	req.Host = rp.origin.Host
	req.URL.Host = rp.origin.Host
	req.URL.Scheme = rp.origin.Scheme
	req.RequestURI = ""
}

func (rp *reloaderProxy) modify(r *http.Response) error {

	for name, interceptor := range rp.interceptors {
		err := interceptor(r)
		if err != nil {
			return fmt.Errorf("interceptor '%s' failed to run: %v", name, err)
		}
	}

	return nil
}

func (rp *reloaderProxy) Run() error {
	err := http.ListenAndServe(rp.endpoint.Host, rp.proxy)
	if err != nil {
		return fmt.Errorf("reloader proxy failed: %v", err)
	}
	return nil
}

type ConfigurerNew interface {
	SetOrigin(url string) error
	SetEndpoint(url string) error
	AddInterceptor(name string, interceptor func(*http.Response) error) error
}

type configurerPoolNew struct {
	pool []func(*reloaderProxy) error
}

func (c *configurerPoolNew) SetOrigin(url string) error {

	o, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse origin url: %v", err)
	}

	c.pool = append(c.pool, func(rp *reloaderProxy) error {
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

	c.pool = append(c.pool, func(rp *reloaderProxy) error {
		rp.endpoint = o
		return nil
	})

	return nil
}

func (c *configurerPoolNew) AddInterceptor(name string, interceptor func(*http.Response) error) error {

	c.pool = append(c.pool, func(rp *reloaderProxy) error {

		if _, ok := rp.interceptors[name]; ok {
			return fmt.Errorf("failed to load an new interceptor: it already exists an interceptor named %s", name)
		}

		rp.interceptors[name] = interceptor

		return nil
	})
	return nil
}

func New(options ...func(ConfigurerNew) error) (*reloaderProxy, error) {

	configurer := &configurerPoolNew{}

	for _, option := range options {
		err := option(configurer)
		if err != nil {
			return nil, fmt.Errorf("failed to load options: %v", err)
		}
	}

	proxy := &reloaderProxy{
		interceptors: make(map[string]func(*http.Response) error),
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
