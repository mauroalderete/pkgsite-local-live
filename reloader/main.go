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

		lr, err := livereload.New()
		if err != nil {
			return fmt.Errorf("failed to load livereload interceptor")
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
