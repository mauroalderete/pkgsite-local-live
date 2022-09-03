// Package cmd implements the command handler
package cmd

import (
	"fmt"
	"log"

	"github.com/mauroalderete/pkgsite-local-live/reloader/reloader-proxy/interceptor/livereload"
	"github.com/mauroalderete/pkgsite-local-live/reloader/reloader-proxy/proxy"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "reloader-proxy",
		Short: "Create a proxy webserver instance that inject a reloader script based websocket",
		Long: `reloader-proxy is a proxy webserver that replicate a server endpoint and 
	inject him a javascript snippet on all html requested.
	It allows that the clients listens a websocket server will expected for a reload signal
	to refresh the webpage loaded in the browsers.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			proxy, err := proxy.New(func(cn proxy.Configurer) error {
				cn.Origin(origin)
				cn.Endpoint(endpoint)

				lr, err := livereload.New(func(cn livereload.Configurer) error {
					err := cn.WebserviceInjectable(snippetFilepath)
					if err != nil {
						return fmt.Errorf("failed to configure a webservice injectable resource: %v", err)
					}

					err = cn.ReloadEndpoint(reloadEndpoint)
					if err != nil {
						return fmt.Errorf("failed to configure the endpoint to the webservice injectable: %v", err)
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

			log.Printf("Start proxy at %s to serve the origin %s\n", origin, endpoint)
			err = proxy.Run()
			if err != nil {
				log.Fatalf("Something went wrong while proxy was running: %v", err)
			}

			return nil
		},
	}

	// store the url to the backend endpoint passed by arguments
	origin string

	// store the url to the frontend endpoint passed by arguments
	endpoint string

	// store the path to the file that contains the snippet to inject by livereload.livereaload interceptor.
	snippetFilepath string

	// store the url to the reload microservice requeried by the livereload.livereaload interceptor.
	reloadEndpoint string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&origin, "origin", "o", "", "URL to endpoint that the proxy must be replicate.")
	rootCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "URL to expose origin modified.")
	rootCmd.Flags().StringVarP(&snippetFilepath, "snippet", "s", "", "filepath that contains the html snippet to inject in all html page requested by clients.")
	rootCmd.Flags().StringVarP(&reloadEndpoint, "reloadendpoint", "r", "", "URL to reload microservice endpoint.")
	rootCmd.MarkFlagRequired("origin")
	rootCmd.MarkFlagRequired("endpoint")
	rootCmd.MarkFlagRequired("snippet")
	rootCmd.MarkFlagRequired("reloadendpoint")
}
