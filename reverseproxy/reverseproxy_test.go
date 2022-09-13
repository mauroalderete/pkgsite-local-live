package reverseproxy

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/mauroalderete/pkgsite-local-live/interceptor"
)

type interceptorFake struct {
	rules   []interceptor.InterceptorRuler
	handler interceptor.InterceptorHandler
}

func (i *interceptorFake) Rules() []interceptor.InterceptorRuler {
	return i.rules
}

func (i *interceptorFake) Handler() interceptor.InterceptorHandler {
	return i.handler
}

func TestNew(t *testing.T) {

	t.Run("missing one", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			return nil
		})
		if err == nil {
			t.Errorf("expect an error, got error nil")
			return
		}
	})

	t.Run("missing two", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			err := c.Origin("loclahost")
			if err != nil {
				return err
			}
			return nil
		})
		if err == nil {
			t.Errorf("expect an error, got error nil")
			return
		}
	})

	t.Run("ok without interceptor", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			err := c.Origin("localhost:8080")
			if err != nil {
				return err
			}
			err = c.Public("localhost:9090")
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			t.Errorf("expect error nil, got '%s'", err)
			return
		}
	})

	t.Run("ok with interceptor", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			err := c.Origin("localhost:8080")
			if err != nil {
				return err
			}
			err = c.Public("localhost:9090")
			if err != nil {
				return err
			}
			err = c.AddInterceptor("fake", &interceptorFake{})
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			t.Errorf("expect error nil, got '%s'", err)
			return
		}
	})

	t.Run("set origin fail", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			err := c.Origin("::")
			if err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			t.Errorf("expect an error, got error nil")
			return
		}
	})

	t.Run("set public fail", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			err := c.Public("::")
			if err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			t.Errorf("expect an error, got error nil")
			return
		}
	})

	t.Run("set interceptor repeated", func(t *testing.T) {
		_, err := New(func(c Configurer) error {

			err := c.AddInterceptor("fake", &interceptorFake{})
			if err != nil {
				return err
			}

			err = c.AddInterceptor("fake", &interceptorFake{})
			if err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			t.Errorf("expect an error, got error nil")
			return
		}
	})
}

func TestDirector(t *testing.T) {
	rp := &ReverseProxy{}
	rp.origin = &url.URL{}
	rp.origin.Host = "localhost"
	rp.origin.Scheme = "http"

	request := &http.Request{}
	request.URL = &url.URL{}

	rp.director(request)
	if request.Host != rp.origin.Host ||
		request.URL.Host != rp.origin.Host ||
		request.URL.Scheme != rp.origin.Scheme ||
		request.RequestURI != "" {
		t.Errorf("expected the same request parameters, got %v %v %v", request.Host, request.URL.Host, request.URL.Scheme)
		return
	}
}

func TestModify(t *testing.T) {
	t.Run("ok", func(t *testing.T) {

		rp := &ReverseProxy{}
		rp.origin = &url.URL{}
		rp.origin.Host = "localhost"
		rp.origin.Scheme = "http"

		i := &interceptorFake{}
		i.rules = make([]interceptor.InterceptorRuler, 0)
		i.handler = func(*http.Response) error {
			return nil
		}

		rp.interceptors = make(map[string]interceptor.Interceptor)
		rp.interceptors["a"] = i

		response := &http.Response{}

		err := rp.modify(response)
		if err != nil {
			t.Errorf("expected error nil, got '%v'", err)
		}
	})

	t.Run("with rule failed", func(t *testing.T) {

		rp := &ReverseProxy{}
		rp.origin = &url.URL{}
		rp.origin.Host = "localhost"
		rp.origin.Scheme = "http"

		i := &interceptorFake{}
		i.rules = make([]interceptor.InterceptorRuler, 0)
		i.rules = append(i.rules, func(r *http.Response) bool { return false })
		i.handler = func(*http.Response) error {
			return nil
		}

		rp.interceptors = make(map[string]interceptor.Interceptor)
		rp.interceptors["a"] = i

		response := &http.Response{}

		err := rp.modify(response)
		if err != nil {
			t.Errorf("expected error nil, got '%v'", err)
		}
	})

	t.Run("ok", func(t *testing.T) {

		rp := &ReverseProxy{}
		rp.origin = &url.URL{}
		rp.origin.Host = "localhost"
		rp.origin.Scheme = "http"

		i := &interceptorFake{}
		i.rules = make([]interceptor.InterceptorRuler, 0)
		i.handler = func(*http.Response) error {
			return fmt.Errorf("some was wrong")
		}

		rp.interceptors = make(map[string]interceptor.Interceptor)
		rp.interceptors["a"] = i

		response := &http.Response{}

		err := rp.modify(response)
		if err == nil {
			t.Errorf("expected an error, got error nil")
		}
	})
}
