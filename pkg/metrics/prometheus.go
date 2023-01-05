package metrics

import (
	"fmt"
	"net/http"
	"tlsmonitor/pkg/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	hostCertExpirationTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "host_cert_expiration_seconds",
			Help: "The time remaning before expiration in seconds",
		},
		[]string{"fqdn", "ip", "certAuthority"},
	)
	fileCertExpirationTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "file_cert_expiration_seconds",
			Help: "The time remaining before expiration in seconds",
		},
		[]string{"filename", "cn", "dnsNames", "certAuthority"},
	)
)

func UpdateHostCert(fqdn string, ip string, ca string, value float64) {
	hostCertExpirationTime.WithLabelValues(fqdn, ip, ca).Set(value)
}

func UpdateFileCert(filename string, cn string, dnsNames string, certAuthority string, value float64) {
	fileCertExpirationTime.WithLabelValues(filename, cn, dnsNames, certAuthority).Set(value)
}

func InitMetrics() {
	prometheus.MustRegister(hostCertExpirationTime)
	prometheus.MustRegister(fileCertExpirationTime)
}

func ServeMetrics(config config.Metrics) {
	if config.Port == 0 {
		config.Port = 9090
	}
	if config.Path == "" {
		config.Path = "/metrics"
	}
	http.Handle(config.Path, promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
