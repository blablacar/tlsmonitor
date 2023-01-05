package cmd

import (
	"os"
	"time"
	"tlsmonitor/pkg/config"
	"tlsmonitor/pkg/tls"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(directCmd)
}

var directCmd = &cobra.Command{
	Use:   "direct",
	Short: "Only run the certificate checks and stop",
	RunE: func(cmd *cobra.Command, args []string) error {
		var serverConfig config.Config

		if Debug {
			log.SetLevel(log.DebugLevel)
		}

		serverConfig, err := LoadConfig()
		if err != nil {
			log.WithError(err).Error("Unable to loadvonfig")
			os.Exit(1)
		}

		err = StartMetrics(&serverConfig)
		if err != nil {
			log.WithError(err).Error("Unable to start metrics")
			os.Exit(1)
		}

		log.Debug(serverConfig)
		tls.StartChecks(&serverConfig)

		if serverConfig.Metrics.WaitScraping {
			time.Sleep(time.Second * 60)
		}

		return nil
	},
}
