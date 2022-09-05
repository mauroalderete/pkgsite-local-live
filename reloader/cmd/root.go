// Package cmd implements the command handler
package cmd

import (
	"fmt"
	"log"

	"github.com/mauroalderete/pkgsite-local-live/reloader/server"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "reloader",
		Short: "Create a proxy webserver instance that inject a reloader script based websocket",
		Long: `reloader is a proxy webserver that replicate a server endpoint and 
	inject him a javascript snippet on all html requested.
	It allows that the clients listens a websocket server will expected for a reload signal
	to refresh the webpage loaded in the browsers.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			srv, err := server.New(func(c server.Configurator) error {
				err := c.Origin(origin)
				if err != nil {
					return fmt.Errorf("failed to configure the origin address to the server instance:%v", err)
				}
				err = c.Public(public)
				if err != nil {
					return fmt.Errorf("failed to configure the public address to the server instance:%v", err)
				}
				err = c.ReloadSnippet(snippetFilepath)
				if err != nil {
					return fmt.Errorf("failed to configure the reload snippet path to the server instance:%v", err)
				}
				return nil
			})
			if err != nil {
				log.Fatalf("Something went wrong to configure the server: %v", err)
			}

			log.Printf("Start server at %s to serve the origin %s\n", origin, public)
			err = srv.Run()
			if err != nil {
				log.Fatalf("Something went wrong while proxy was running: %v", err)
			}

			return nil
		},
	}

	// store the url to the backend endpoint passed by arguments
	origin string

	// store the url to the frontend public passed by arguments
	public string

	// store the path to the file that contains the snippet to inject by livereload.livereaload interceptor.
	snippetFilepath string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&origin, "origin", "o", "", "URL to endpoint that the proxy must be replicate.")
	rootCmd.Flags().StringVarP(&public, "public", "p", "", "URL to expose origin modified.")
	rootCmd.Flags().StringVarP(&snippetFilepath, "snippet", "s", "", "filepath that contains the html snippet to inject in all html page requested by clients.")
	rootCmd.MarkFlagRequired("origin")
	rootCmd.MarkFlagRequired("public")
	rootCmd.MarkFlagRequired("snippet")
}
