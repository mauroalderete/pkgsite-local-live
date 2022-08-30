package main

import (
	"fmt"
	"log"

	"github.com/mauroalderete/pkgsite-local-live/reloader/interceptor/livereload"
	"github.com/mauroalderete/pkgsite-local-live/reloader/reloaderproxy"
)

func main() {

	proxy, err := reloaderproxy.New(func(cn reloaderproxy.ConfigurerNew) error {
		cn.SetOrigin("http://localhost:3000")
		cn.SetEndpoint("http://localhost:9090")

		lr, err := livereload.New(func(cn livereload.ConfigurerNew) error {
			err := cn.WebserviceInjectable("./interceptor/livereload/websocket.html")
			if err != nil {
				return fmt.Errorf("failed to configure a webservice injectable resource: %v", err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to load livereload interceptor: %v", err)
		}
		cn.AddInterceptor("livereload", lr)

		return nil
	})

	if err != nil {
		log.Fatalf("Something went wrong to configure the proxy: %v", err)
	}

	err = proxy.Run()
	if err != nil {
		log.Fatalf("Something went wrong while proxy was running: %v", err)
	}
}
