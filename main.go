package main

import (
	"log"

	"github.com/mauroalderete/pkgsite-local-live/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
}
