package cmd

import (
	"tlsmonitor/pkg/config"
	"tlsmonitor/pkg/metrics"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() (config.Config, error) {
	var serverConfig config.Config

	viper.SetConfigFile(ConfigPath)
	log.Debug("Config path; %s", ConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Error("Cannot read: ", ConfigPath)
		return serverConfig, err
	}

	err := viper.Unmarshal(&serverConfig)
	if err != nil {
		log.WithError(err).Error("unable to decode into struct")
		return serverConfig, err
	}

	log.Debug(serverConfig)
	return serverConfig, nil
}

func StartMetrics(serverConfig *config.Config) error {
	if serverConfig.Metrics.Enabled {
		metrics.InitMetrics()
		go func() {
			metrics.ServeMetrics(serverConfig.Metrics)
		}()
	}
	return nil
}
