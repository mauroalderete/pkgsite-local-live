package websocketconnections

import (
	"net/http"
	"testing"
)

type responseWriterFacke struct {
	header      http.Header
	write       func([]byte) (int, error)
	writeHeader func(statusCode int)
}

func (r *responseWriterFacke) Header() http.Header {
	return r.header
}

func (r *responseWriterFacke) Write(data []byte) (int, error) {
	return r.write(data)
}

func (r *responseWriterFacke) WriteHeader(statusCode int) {
	r.writeHeader(statusCode)
}

func TestNew(t *testing.T) {

	t.Run("", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			return nil
		})
		if err == nil {
			t.Errorf("expected an error, got error nil")
		}
	})

	t.Run("", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			return c.Request(&http.Request{})
		})
		if err == nil {
			t.Errorf("expected an error, got error nil")
		}
	})

	t.Run("", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			return c.ResponseWriter(&responseWriterFacke{})
		})
		if err == nil {
			t.Errorf("expected an error, got error nil")
		}
	})

	t.Run("", func(t *testing.T) {
		_, err := New(func(c Configurer) error {
			err := c.Request(&http.Request{})
			if err != nil {
				return err
			}

			err = c.ResponseWriter(&responseWriterFacke{})
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			t.Errorf("expected error nil, got '%v'", err)
		}
	})
}
