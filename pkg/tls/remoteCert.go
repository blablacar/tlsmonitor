package tls

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
	"tlsmonitor/pkg/config"
	"tlsmonitor/pkg/metrics"

	log "github.com/sirupsen/logrus"
)

func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleep)
		log.WithFields(log.Fields{"retry_count": i}).Debug("Retrying after error:", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func CheckHostCert(target string, config *config.Config) {
	log.WithFields(log.Fields{"target": target}).Info("Checking cert")
	targetPort := strings.Split(target, ":")
	port := "443"
	fqdn := targetPort[0]
	if len(targetPort) > 2 {
		log.WithFields(log.Fields{"target": target}).Error("Host value not valid")
		return
	} else if len(targetPort) == 2 {
		port = targetPort[1]
	}

	// Set to insecure because of internal PKI certiticates or self signed wrong hosts ...
	dialConfig := tls.Config{ServerName: fqdn, InsecureSkipVerify: true}

	var ips []net.IP
	errLookup := retry(config.NbDialRetry, 5*time.Second, func() (err error) {
		ips, err = net.LookupIP(fqdn)
		return
	})
	if errLookup != nil {
		log.WithFields(log.Fields{"fqdn": fqdn}).Error("Lookup IP failed")
		if config.Metrics.Enabled {
			metrics.UpdateHostCert(fqdn, "0.0.0.0", "NA", 0)
		}
		return
	}

	timeNow := time.Now()
	alertYears := 0
	alertMonths := 0
	alertDays := 40

	for _, ip := range ips {
		if ip.To4() != nil {
			var conn *tls.Conn
			errDial := retry(config.NbDialRetry, 5*time.Second, func() (err error) {
				conn, err = tls.Dial("tcp", fmt.Sprint(ip, ":", port), &dialConfig)
				return
			})
			if errDial != nil {
				log.WithFields(log.Fields{"status": "FAILED", "target": target, "ip": ip}).Warn(errDial.Error())
				if config.Metrics.Enabled {
					metrics.UpdateHostCert(fqdn, ip.String(), "NA", 0)
				}
				return
			}
			defer conn.Close()
			// Certificate is first in chain
			endCert := conn.ConnectionState().PeerCertificates[0]

			remainingTime := endCert.NotAfter.Unix() - timeNow.Unix()
			if timeNow.After(endCert.NotAfter) {
				log.WithFields(log.Fields{
					"status":        "EXPIRED",
					"fqdn":          fqdn,
					"ip":            ip,
					"expiration":    endCert.NotAfter,
					"remainingtime": remainingTime,
				}).Warn("Cert is expired: ", endCert.Subject.CommonName)
			} else if timeNow.AddDate(alertYears, alertMonths, alertDays).After(endCert.NotAfter) {
				log.WithFields(log.Fields{
					"status":        "SOON",
					"fqdn":          fqdn,
					"ip":            ip,
					"expiration":    endCert.NotAfter,
					"remainingtime": remainingTime,
				}).Warn("Cert will expired soon: ", endCert.Subject.CommonName)
			} else {
				log.WithFields(log.Fields{
					"status":        "VALID",
					"fqdn":          fqdn,
					"ip":            ip,
					"expiration":    endCert.NotAfter,
					"remainingtime": remainingTime,
				}).Info("Cert is valid: ", endCert.Subject.CommonName)
			}
			if config.Metrics.Enabled {
				metrics.UpdateHostCert(fqdn, ip.String(), endCert.Issuer.CommonName, float64(remainingTime))
			}
		}
	}
}
