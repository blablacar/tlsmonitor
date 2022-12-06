package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tlsmonitor",
	Short: "Monitoring certificate expiration",
	Long: `TLS Monitor will parse certificate to get the time remaining before
					expiration, it can expose metrics in prometheus format.`,
}

var (
	Debug      bool
	ConfigPath string
)

func Execute() {
	log.SetFormatter(&log.JSONFormatter{})
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Enable debug logs")
	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "./config.yaml", "Config file path")

	if Debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error("an error occurred")
		os.Exit(1)
	}
}
