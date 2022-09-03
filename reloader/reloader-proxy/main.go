package main

import (
	"log"

	"github.com/mauroalderete/pkgsite-local-live/reloader/reloader-proxy/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
}
