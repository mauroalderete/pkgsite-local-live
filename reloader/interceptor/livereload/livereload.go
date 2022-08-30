package livereload

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mauroalderete/pkgsite-local-live/reloader/interceptor"
)

var javascriptToInject string

type Livereload struct {
	rules   []interceptor.InterceptorRuler
	handler interceptor.InterceptorHandler
}

func (l *Livereload) Rules() []interceptor.InterceptorRuler {
	return l.rules
}

func (l *Livereload) Handler() interceptor.InterceptorHandler {
	return l.handler
}

func statusCodeRule(r *http.Response) bool {
	return r.StatusCode == 200
}

func contentTypeRule(r *http.Response) bool {
	isTextHML := false
	for _, v := range r.Header["Content-Type"] {
		if strings.Contains(v, "text/html") {
			isTextHML = true
			break
		}
	}

	return isTextHML
}

func getBody(r *http.Response) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the body: %v", err)
	}
	err = r.Body.Close()
	if err != nil {
		return "", fmt.Errorf("failed to terminate the body: %v", err)
	}

	return string(body), nil
}

func handler(r *http.Response) error {
	body, err := getBody(r)
	if err != nil {
		return fmt.Errorf("failed to get body: %v", err)
	}

	r.Body = io.NopCloser(strings.NewReader("Hola World" + body))

	return nil
}

func New() (interceptor.Interceptor, error) {

	livereload := &Livereload{
		rules:   []interceptor.InterceptorRuler{statusCodeRule, contentTypeRule},
		handler: handler,
	}

	return livereload, nil
}
