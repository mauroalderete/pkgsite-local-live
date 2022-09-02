package livereload

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mauroalderete/pkgsite-local-live/reloader/reloaderproxy/interceptor"

	"regexp"
	"strings"
)

type Livereload struct {
	webserviceInjectable string
	rules                []interceptor.InterceptorRuler
}

func (l *Livereload) Rules() []interceptor.InterceptorRuler {
	return l.rules
}

func (l *Livereload) Handler() interceptor.InterceptorHandler {
	return func(r *http.Response) error {
		content, err := getBody(r)
		if err != nil {
			return fmt.Errorf("failed to get body: %v", err)
		}

		exp := regexp.MustCompile("</body>")
		location := exp.FindIndex([]byte(content))

		contentModified := content[:location[0]-1]
		contentModified += fmt.Sprintf("\n%s\n", l.webserviceInjectable)
		contentModified += content[location[0]:]

		r.Body = io.NopCloser(strings.NewReader(contentModified))

		return nil
	}
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

func hasOneBodyTagRule(r *http.Response) bool {
	content, err := getBody(r)
	if err != nil {
		return false
	}

	openTagExp := regexp.MustCompile("<body>")
	closeTagExp := regexp.MustCompile("</body>")
	openTagMatchs := openTagExp.FindAll([]byte(content), -1)
	closeTagMatchs := closeTagExp.FindAll([]byte(content), -1)

	r.Body = io.NopCloser(strings.NewReader(content))

	return len(openTagMatchs) == 1 && len(closeTagMatchs) == 1
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

type ConfigurerNew interface {
	WebserviceInjectable(path string) error
}

type configurerNew struct {
	pool []func(l *Livereload) error
}

func (c *configurerNew) WebserviceInjectable(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	c.pool = append(c.pool, func(l *Livereload) error {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed load webservice injectable resource from %s: %v", path, err)
		}
		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to access at the content of webservice injectable resource: %v", err)
		}

		l.webserviceInjectable = string(content)
		return nil
	})

	return nil
}

func New(options ...func(ConfigurerNew) error) (interceptor.Interceptor, error) {

	livereload := &Livereload{
		rules: []interceptor.InterceptorRuler{
			statusCodeRule,
			contentTypeRule,
			hasOneBodyTagRule,
		},
	}

	configurer := &configurerNew{}

	for _, option := range options {
		err := option(configurer)
		if err != nil {
			return nil, fmt.Errorf("failed to load the configuration: %v", err)
		}
	}

	for _, config := range configurer.pool {
		err := config(livereload)
		if err != nil {
			return nil, fmt.Errorf("failed to apply the configuration: %v", err)
		}
	}

	if len(livereload.webserviceInjectable) == 0 {
		return nil, fmt.Errorf("a webserviceInjectable is required")
	}

	return livereload, nil
}
