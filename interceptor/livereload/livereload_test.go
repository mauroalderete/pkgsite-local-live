package livereload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestConfigureOpenFileNil(t *testing.T) {
	config := &configurer{}

	err := config.OpenFile(nil)
	if err == nil {
		t.Errorf("want an error, got nil")
	}
}

func TestConfigureReadAllNil(t *testing.T) {
	config := &configurer{}

	err := config.ReadAll(nil)
	if err == nil {
		t.Errorf("want an error, got nil")
	}
}

func TestConfigureUpgradeEndpointEmpty(t *testing.T) {
	config := &configurer{}

	err := config.UpgradeEndpoint("")
	if err == nil {
		t.Errorf("want an error, got nil")
	}
}

func TestConfigureUpgradeEndpoint(t *testing.T) {

	livereloadMock := &Livereload{}
	config := &configurer{}

	err := config.UpgradeEndpoint("some address")
	if err != nil {
		t.Errorf("want error nil, got '%s'", err)
	}

	err = config.pool[0](livereloadMock)
	if err != nil {
		t.Errorf("want error nil, got '%s'", err)
		return
	}
}

func TestWebserviceInjectableEmpty(t *testing.T) {
	conf := &configurer{}

	err := conf.WebserviceInjectable("")
	if err == nil {
		t.Errorf("got error nil, want an error")
		return
	}

	t.Run("path wrong", func(t *testing.T) {
		err := conf.WebserviceInjectable("")
		if err == nil {
			t.Errorf("want an error, got error nil")
			return
		}
	})
}

func TestWebserviceInjectableWrongPath(t *testing.T) {
	conf := &configurer{}

	conf.OpenFile(func(name string) (*os.File, error) {
		return nil, fmt.Errorf("file not found")
	})

	err := conf.WebserviceInjectable("file")
	if err != nil {
		t.Errorf("want error nil, got %s", err)
		return
	}

	livereloadMock := &Livereload{}

	err = conf.pool[0](livereloadMock)
	if err != nil {
		t.Errorf("want error nil, got error %v", err)
		return
	}

	err = conf.pool[1](livereloadMock)
	if err == nil {
		t.Errorf("want an error, got error nil")
		return
	}
}

func TestWebserviceInjectableReadFailed(t *testing.T) {

	// Prepares mocks
	livereloadMock := &Livereload{}
	conf := &configurer{}

	conf.OpenFile(func(name string) (*os.File, error) {
		return &os.File{}, nil
	})

	conf.ReadAll(func(r io.Reader) ([]byte, error) {
		return []byte{}, fmt.Errorf("failed to read all")
	})

	for _, config := range conf.pool {
		err := config(livereloadMock)
		if err != nil {
			t.Errorf("want error nil, got 'failed to apply the configuration: %v'", err)
			return
		}
	}
	conf.pool = conf.pool[:0]

	// Test core
	err := conf.WebserviceInjectable("file")
	if err != nil {
		t.Errorf("want error nil, got %s", err)
		return
	}

	err = conf.pool[0](livereloadMock)
	if err == nil {
		t.Errorf("want an error, got error nil")
		return
	}
}

func TestWebserviceInjectableEmptyContent(t *testing.T) {
	// Prepares mocks
	livereloadMock := &Livereload{}
	conf := &configurer{}

	conf.OpenFile(func(name string) (*os.File, error) {
		return &os.File{}, nil
	})

	conf.ReadAll(func(r io.Reader) ([]byte, error) {
		return []byte{}, nil
	})

	for _, config := range conf.pool {
		err := config(livereloadMock)
		if err != nil {
			t.Errorf("want error nil, got 'failed to apply the configuration: %v'", err)
			return
		}
	}
	conf.pool = conf.pool[:0]

	// Test core
	err := conf.WebserviceInjectable("file")
	if err != nil {
		t.Errorf("want error nil, got %s", err)
		return
	}

	err = conf.pool[0](livereloadMock)
	if err != nil {
		t.Errorf("want error nil, got %s", err)
		return
	}

	const expected = ""

	if livereloadMock.webserviceInjectable != expected {
		t.Errorf("want '%s', got '%s'", expected, livereloadMock.webserviceInjectable)
		return
	}
}

func TestWebserviceInjectableSuccefful(t *testing.T) {
	// Prepares mocks
	livereloadMock := &Livereload{}
	conf := &configurer{}

	conf.OpenFile(func(name string) (*os.File, error) {
		return &os.File{}, nil
	})

	conf.ReadAll(func(r io.Reader) ([]byte, error) {
		return []byte("some content"), nil
	})

	for _, config := range conf.pool {
		err := config(livereloadMock)
		if err != nil {
			t.Errorf("want error nil, got 'failed to apply the configuration: %v'", err)
			return
		}
	}
	conf.pool = conf.pool[:0]

	// Test core
	err := conf.WebserviceInjectable("file")
	if err != nil {
		t.Errorf("want error nil, got %s", err)
		return
	}

	err = conf.pool[0](livereloadMock)
	if err != nil {
		t.Errorf("want error nil, got %s", err)
		return
	}

	const expected = "some content"

	if livereloadMock.webserviceInjectable != expected {
		t.Errorf("want '%s', got '%s'", expected, livereloadMock.webserviceInjectable)
		return
	}
}

func TestRules(t *testing.T) {
	livereload, err := New(
		func(c Configurer) error {
			c.OpenFile(func(name string) (*os.File, error) {
				return &os.File{}, nil
			})

			c.ReadAll(func(r io.Reader) ([]byte, error) {
				return []byte("some content"), nil
			})
			return nil
		},
		func(c Configurer) error {
			err := c.UpgradeEndpoint("address")
			if err != nil {
				return err
			}

			err = c.WebserviceInjectable("some pathfile")
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		t.Errorf("failed to instance livereload mock, got '%v'", err)
		return
	}

	expected := 3
	got := len(livereload.Rules())
	if expected != got {
		t.Errorf("expected %d rules, got %d", expected, got)
		return
	}
}

func TestNew(t *testing.T) {

	t.Run("without openFile", func(t *testing.T) {
		_, err := New(
			func(c Configurer) error {
				c.ReadAll(func(r io.Reader) ([]byte, error) {
					return []byte("some content"), nil
				})
				return nil
			},
			func(c Configurer) error {
				err := c.UpgradeEndpoint("address")
				if err != nil {
					return err
				}

				err = c.WebserviceInjectable("some pathfile")
				if err != nil {
					return err
				}

				return nil
			},
		)
		if err == nil {
			t.Errorf("want an error, got nil")
			return
		}
	})

	t.Run("without readAll", func(t *testing.T) {
		_, err := New(
			func(c Configurer) error {
				c.OpenFile(func(name string) (*os.File, error) {
					return &os.File{}, nil
				})
				return nil
			},
			func(c Configurer) error {
				err := c.UpgradeEndpoint("address")
				if err != nil {
					return err
				}

				err = c.WebserviceInjectable("some pathfile")
				if err != nil {
					return err
				}

				return nil
			},
		)
		if err == nil {
			t.Errorf("want an error, got nil")
			return
		}
	})

	t.Run("without webserviceInjectable", func(t *testing.T) {
		_, err := New(
			func(c Configurer) error {
				c.OpenFile(func(name string) (*os.File, error) {
					return &os.File{}, nil
				})

				c.ReadAll(func(r io.Reader) ([]byte, error) {
					return []byte("some content"), nil
				})
				return nil
			},
			func(c Configurer) error {
				err := c.UpgradeEndpoint("address")
				if err != nil {
					return err
				}

				return nil
			},
		)
		if err == nil {
			t.Errorf("want an error, got nil")
			return
		}
	})

	t.Run("without endpoint", func(t *testing.T) {
		_, err := New(
			func(c Configurer) error {
				c.OpenFile(func(name string) (*os.File, error) {
					return &os.File{}, nil
				})

				c.ReadAll(func(r io.Reader) ([]byte, error) {
					return []byte("some content"), nil
				})
				return nil
			},
			func(c Configurer) error {
				err := c.WebserviceInjectable("some pathfile")
				if err != nil {
					return err
				}

				return nil
			},
		)
		if err == nil {
			t.Errorf("want an error, got nil")
			return
		}
	})

	t.Run("fail prepare some config", func(t *testing.T) {
		_, err := New(
			func(c Configurer) error {
				return fmt.Errorf("some was wrong")
			},
		)
		if err == nil {
			t.Errorf("want an error, got nil")
			return
		}
	})

	t.Run("fail apply some config", func(t *testing.T) {
		_, err := New(
			func(c Configurer) error {
				c.OpenFile(func(name string) (*os.File, error) {
					return &os.File{}, fmt.Errorf("some was wrong")
				})

				c.ReadAll(func(r io.Reader) ([]byte, error) {
					return []byte("some content"), nil
				})
				return nil
			},
			func(c Configurer) error {
				err := c.WebserviceInjectable("some pathfile")
				if err != nil {
					return err
				}

				return nil
			},
		)
		if err == nil {
			t.Errorf("want an error, got nil")
			return
		}
	})
}

func TestStatusCodeRule(t *testing.T) {
	cases := map[string]struct {
		response *http.Response
		expected bool
	}{
		"200": {&http.Response{StatusCode: 200}, true},
		"304": {&http.Response{StatusCode: 304}, true},
		"404": {&http.Response{StatusCode: 404}, false},
		"501": {&http.Response{StatusCode: 501}, false},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			if statusCodeRule(c.response) != c.expected {
				t.Errorf("expected %v, got %v", c.expected, !c.expected)
				return
			}
		})
	}
}

func TestContentTypeRule(t *testing.T) {
	cases := map[string]struct {
		response *http.Response
		expected bool
	}{
		"html":    {&http.Response{Header: map[string][]string{"Content-Type": {"text/html"}}}, true},
		"json":    {&http.Response{Header: map[string][]string{"Content-Type": {"text/json"}}}, false},
		"image":   {&http.Response{Header: map[string][]string{"Content-Type": {"image/*"}}}, false},
		"unknown": {&http.Response{Header: map[string][]string{"Content-Type": {"some other"}}}, false},
		"without": {&http.Response{}, true},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			if contentTypeRule(c.response) != c.expected {
				t.Errorf("expected %v, got %v", c.expected, !c.expected)
				return
			}
		})
	}
}

type nopCloserMock struct {
	io.Reader
	close func() error
}

func (ncm *nopCloserMock) Close() error { return ncm.close() }

func TestHasOneBodyTagRule(t *testing.T) {
	cases := map[string]struct {
		response *http.Response
		reader   func(r io.Reader) ([]byte, error)
		expected bool
	}{
		"empty":              {&http.Response{Body: io.NopCloser(strings.NewReader(""))}, func(r io.Reader) ([]byte, error) { return io.ReadAll(r) }, false},
		"only open tag":      {&http.Response{Body: io.NopCloser(strings.NewReader("abc<body>abc"))}, func(r io.Reader) ([]byte, error) { return io.ReadAll(r) }, false},
		"only close tag":     {&http.Response{Body: io.NopCloser(strings.NewReader("abc</body>abc"))}, func(r io.Reader) ([]byte, error) { return io.ReadAll(r) }, false},
		"two open one close": {&http.Response{Body: io.NopCloser(strings.NewReader("abc<body>abc<body>abc</body>abc"))}, func(r io.Reader) ([]byte, error) { return io.ReadAll(r) }, false},
		"one open two close": {&http.Response{Body: io.NopCloser(strings.NewReader("abc<body>abc</body>abc</body>abc"))}, func(r io.Reader) ([]byte, error) { return io.ReadAll(r) }, false},
		"reader failed":      {&http.Response{Body: io.NopCloser(strings.NewReader(""))}, func(r io.Reader) ([]byte, error) { return nil, fmt.Errorf("some was wrong") }, false},
		"close filed":        {&http.Response{Body: &nopCloserMock{strings.NewReader(""), func() error { return fmt.Errorf("some was wrong") }}}, func(r io.Reader) ([]byte, error) { return io.ReadAll(r) }, false},
		"succefull":          {&http.Response{Body: io.NopCloser(strings.NewReader("abc<body>abc</body>abc"))}, func(r io.Reader) ([]byte, error) { return io.ReadAll(r) }, true},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			if hasOneBodyTagRule(c.response, c.reader) != c.expected {
				t.Errorf("expected %v, got %v", c.expected, !c.expected)
				return
			}
		})
	}
}

func TestInterceptor(t *testing.T) {
	livereload := &Livereload{}

	livereload.webserviceInjectable = "1234567890"
	livereload.readAll = io.ReadAll

	response := &http.Response{}
	response.Body = io.NopCloser(strings.NewReader("abc<body>abc</body>abc"))
	response.Header = make(http.Header)

	interceptorCallback := livereload.Handler()
	err := interceptorCallback(response)
	if err != nil {
		t.Errorf("failed to try execute the handler: %v", err)
		return
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to try read response modified: %v", err)
		return
	}

	expected := fmt.Sprintf("abc<body>abc\n%s\n</body>abc", livereload.webserviceInjectable)

	if string(result) != expected {
		t.Errorf("expected '%s', got '%s'", expected, string(result))
		return
	}
}

func TestInterceptorError(t *testing.T) {
	livereload := &Livereload{}

	livereload.webserviceInjectable = "1234567890"
	livereload.readAll = func(r io.Reader) ([]byte, error) { return []byte{}, fmt.Errorf("some was wrong") }

	response := &http.Response{}
	response.Body = io.NopCloser(strings.NewReader("abc<body>abc</body>abc"))
	response.Header = make(http.Header)

	interceptorCallback := livereload.Handler()
	err := interceptorCallback(response)
	if err == nil {
		t.Errorf("excpected an error, but error is nil")
		return
	}
}
