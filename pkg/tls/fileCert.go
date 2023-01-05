package tls

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"strings"
	"time"
	"tlsmonitor/pkg/config"
	"tlsmonitor/pkg/metrics"

	log "github.com/sirupsen/logrus"
)

func CheckFileCert(certPath string, config *config.Config) {
	f, errReadFile := ioutil.ReadFile(certPath)
	if errReadFile != nil {
		log.WithFields(log.Fields{"certPath": certPath}).Error(errReadFile.Error())
		return
	}
	pem, _ := pem.Decode(f)
	if pem == nil {
		log.WithFields(log.Fields{"certPath": certPath}).Error("Error decoding pem: ", certPath)
		return
	}
	certs, err := x509.ParseCertificates(pem.Bytes)
	if err != nil {
		log.WithFields(log.Fields{"certPath": certPath}).Error(err.Error())
		return
	}
	log.WithFields(log.Fields{"cn": certs[0].Subject.CommonName, "expiration": certs[0].NotAfter}).
		Info("Checking certificate: ", certs[0].Subject.CommonName)

	remainingTime := certs[0].NotAfter.Unix() - time.Now().Unix()
	if config.Metrics.Enabled {
		metrics.UpdateFileCert(
			certPath,
			certs[0].Subject.CommonName,
			strings.Join(certs[0].DNSNames, ", "),
			certs[0].Issuer.CommonName,
			float64(remainingTime),
		)
	}
}
