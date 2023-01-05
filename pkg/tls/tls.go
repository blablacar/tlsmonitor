package tls

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"tlsmonitor/pkg/config"

	log "github.com/sirupsen/logrus"
)

func StartChecks(config *config.Config) {
	if config.NbDialRetry == 0 {
		config.NbDialRetry = 3
	}

	var wgHosts sync.WaitGroup
	wgHosts.Add(len(config.Hosts))
	for _, host := range config.Hosts {
		go func(host string) {
			CheckHostCert(host, config)
			wgHosts.Done()
		}(host)
	}

	for _, certPath := range config.CertsPaths {
		fm, err := os.Stat(certPath)
		if err != nil {
			log.WithFields(log.Fields{"certPath": certPath}).Error(err.Error())
			continue
		}
		if fm.Mode().IsRegular() {
			go func(certPath string) {
				CheckFileCert(certPath, config)
			}(certPath)
		} else if fm.Mode().IsDir() {
			var files []string
			err := filepath.Walk(certPath, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
					log.Debug("Hidden dir ignored: ", info.Name())
					return filepath.SkipDir
				}
				if strings.HasPrefix(info.Name(), ".") {
					log.Debug("Hidden file ignored: ", info.Name())
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if strings.Contains(info.Name(), "key") {
					log.Debug("Cert checking excluded key file: ", info.Name())
					return nil
				}
				files = append(files, path)
				log.Debug("Cert checking add ", info.Name(), " to list of certs")
				return nil
			})
			if err != nil {
				log.Error(err.Error())
				continue
			}
			for _, file := range files {
				CheckFileCert(file, config)
			}
		}
	}

	wgHosts.Wait()
}
