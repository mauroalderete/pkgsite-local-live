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

			proxy, err := proxy.New(func(cn proxy.ConfigurerNew) error {
				cn.SetOrigin(origin)
				cn.SetEndpoint(endpoint)

				lr, err := livereload.New(func(cn livereload.Configurer) error {
					err := cn.WebserviceInjectable(snippetFilepath)
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

			log.Printf("Start proxy at %s to serve the origin %s\n", origin, endpoint)
			err = proxy.Run()
			if err != nil {
				log.Fatalf("Something went wrong while proxy was running: %v", err)
			}

			return nil
		},
	}
	origin          string
	endpoint        string
	snippetFilepath string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&origin, "origin", "o", "", "URL to endpoint that the proxy must be replicate.")
	rootCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "URL to expose origin modified.")
	rootCmd.Flags().StringVarP(&snippetFilepath, "snippet", "s", "", "filepath that contains the html snippet to inject in all html page requested by clients.")
	rootCmd.MarkFlagRequired("origin")
	rootCmd.MarkFlagRequired("endpoint")
	rootCmd.MarkFlagRequired("snippet")
}