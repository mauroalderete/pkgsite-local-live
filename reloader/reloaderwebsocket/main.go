package main

import (
	"fmt"
	"log"

	"github.com/mauroalderete/pkgsite-local-live/reloader/reloaderwebsocket/server"
)

func main() {
	websocket, err := server.New(func(cn server.ConfigurerNew) error {
		err := cn.Endpoint("localhost:9090")
		if err != nil {
			return fmt.Errorf("failed to configure endpoint of the new reloader websocket: %v", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("failed to create a new teload websocket %v", err)
	}

	websocket.Run()
}
