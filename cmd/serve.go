package cmd

import (
	"log"

	"github.com/Hamaiz/go-rest-eg/config"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := serve.NewServer()
		if err != nil {
			log.Fatal(err)
		}
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
