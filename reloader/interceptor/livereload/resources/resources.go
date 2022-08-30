package resources

import (
	"encoding/base64"
	"log"
)

var webserviceLinked string

var Webservice struct {
	Encoded string
	Decoded string
}

func init() {
	decoded, err := base64.StdEncoding.DecodeString(webserviceLinked)
	if err != nil {
		log.Fatalf("failed to load webservice injectable resource required to livereload interceptor")
	}

	Webservice.Encoded = webserviceLinked
	Webservice.Decoded = string(decoded)
}
