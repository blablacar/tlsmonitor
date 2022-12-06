package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"tlsmonitor/pkg/config"
	"tlsmonitor/pkg/tls"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start monitor and metrics endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		var serverConfig config.Config

		if Debug {
			log.SetLevel(log.DebugLevel)
		}

		serverConfig, err := LoadConfig()
		if err != nil {
			log.WithError(err).Error("Unable to load config")
			os.Exit(1)
		}

		err = StartMetrics(&serverConfig)
		if err != nil {
			log.WithError(err).Error("Unable to start metrics")
			os.Exit(1)
		}

		log.Debug(serverConfig)

		// Start Health check for kubernetes integration
		go func() {
			http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "200 - healthy")
			},
			)
			http.ListenAndServe(":8080", nil)
		}()

		tls.StartChecks(&serverConfig) // Run it a first time before tick
		for range time.Tick(time.Hour * time.Duration(serverConfig.ChecksFrequency)) {
			tls.StartChecks(&serverConfig)
		}

		return nil
	},
}
