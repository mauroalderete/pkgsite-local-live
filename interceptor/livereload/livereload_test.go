package livereload

import (
	"fmt"
	"io"
	"os"
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
		fmt.Printf("asdasd")
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
