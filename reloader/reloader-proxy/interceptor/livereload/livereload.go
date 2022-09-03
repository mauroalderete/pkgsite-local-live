// Package liverreload injects in a html page requested a snippet from a filepath
// to handle a livereload system in client side.
package livereload

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/mauroalderete/pkgsite-local-live/reloader/reloader-proxy/interceptor"

	"regexp"
	"strings"
)

// livereload implements interceptor.Interceptor interface
type livereload struct {
	webserviceInjectable string
	rules                []interceptor.InterceptorRuler
	reloadEndpoint       string
}

// Rules implements interceptor.Interceptor.Rules method.
// Returns a list of interceptor.InterceptorRuler loaded with the rules needed to inject the snippet.
// The rules are loaded during the build of a instance of livereload.livereload.
func (l *livereload) Rules() []interceptor.InterceptorRuler {
	return l.rules
}

// Handler implements interceptor.Interceptor.Handler method.
// Returns a interceptor.InterceptorHandler callback.
//
// The method returned access to the content and inject before tha tag `</body>`
// the snippet passed as option during the build of a instance of liverreload.livereload.
func (l *livereload) Handler() interceptor.InterceptorHandler {
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

// statusCodeRule validates that the response requested has a status code 200.
func statusCodeRule(r *http.Response) bool {
	return r.StatusCode == 200
}

// contentTypeRule validates that the content-type of the response requested is a `text/hmlt`
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

// hasOneBodyTagRule validates that the response requested has only one body HTML tag pair.
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

// getBody allows access to a copy of the body content
// while maintaining open the body in response requested to future readings.
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

// Configurer define the configurable options to build a new instance of livereload.livereload
type Configurer interface {

	// WebserviceInjectable receives the path of file that contains the snippet that must be injected in the body content.
	//
	// Returns an error if failed to get the file or parse it.
	WebserviceInjectable(path string) error

	// ReloadEndpoint set the reload microservice endpoint that of the snippet must be listen to establish the connection with a websocket.
	ReloadEndpoint(url string) error
}

// configurer implement the livereload.Configurer interface.
//
// It stores in a pool the callbacks with the configurable options
// that must be called by the constructor of livereload.livereload to apply the configurations.
type configurer struct {
	pool []func(l *livereload) error
}

// WebserviceInjectable implements livereload.Configurer.WebserviceInjectable method.
//
// Add to the pool the function needed to open and read the snippet file.
func (c *configurer) WebserviceInjectable(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	c.pool = append(c.pool, func(l *livereload) error {
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

func (c *configurer) ReloadEndpoint(url string) error {

	if len(url) == 0 {
		return fmt.Errorf("reload endpoint cannot be empty")
	}

	c.pool = append(c.pool, func(l *livereload) error {
		l.reloadEndpoint = url
		return nil
	})

	return nil
}

// New returns a livereload.livereload instance that implements the interceptor.Interceptor interface.
//
// Receive a list of configurations callback to apply the options.
// It function try to access to the file with the snippet to inject
// and configures the rules needed to identify the request that must be injected.
func New(options ...func(Configurer) error) (interceptor.Interceptor, error) {

	livereload := &livereload{
		rules: []interceptor.InterceptorRuler{
			statusCodeRule,
			contentTypeRule,
			hasOneBodyTagRule,
		},
	}

	configurer := &configurer{}

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

	if len(livereload.reloadEndpoint) == 0 {
		return nil, fmt.Errorf("a reload endpoint is required")
	}

	// Replace the actions in reload snippet
	payload := struct {
		ReloadEndpoint string
	}{
		ReloadEndpoint: livereload.reloadEndpoint,
	}

	tmp := template.New("reload_snippet")
	tmp, _ = tmp.Parse(livereload.webserviceInjectable)

	wr := &bytes.Buffer{}
	err := tmp.Execute(wr, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the injectable snippet script: %v", err)
	}

	livereload.webserviceInjectable = wr.String()

	return livereload, nil
}
