/*
Copyright © 2024 Sundeep Chand
*/
package cmd

import (
	"errors"
	"log"
	"net/http"

	"github.com/SundeepChand/http-proxy/config"
	"github.com/SundeepChand/http-proxy/core"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the proxy to listen for any incoming requests and forward",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := "./examples/basic-http-api/example-config.yml"
		config, err := config.Load(configPath)
		if err != nil {
			log.Fatal("unable to load config", err)
		}

		proxy := core.NewProxyServer(config)

		err = proxy.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server shutting down...")
		} else if err != nil {
			log.Fatal("err in starting server", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
