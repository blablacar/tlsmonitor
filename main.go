package main

import (
	"tlsmonitor/cmd"

	log "github.com/sirupsen/logrus"
	_ "go.uber.org/automaxprocs"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	cmd.Execute()
}
