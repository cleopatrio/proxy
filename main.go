package main

import (
	"os"

	"github.com/cleopatrio/proxy/logger"
	"github.com/cleopatrio/proxy/proxy"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	file, err := os.ReadFile("Proxyfile")
	if err != nil {
		logger.Logger.Warn("Unable to load Proxyfile. Enabled default configuration.")
		os.Exit(1)
	}

	if err := yaml.Unmarshal(file, &proxy.PxFile); err != nil {
		logger.Logger.Fatal("Invalid Proxyfile ", err)
		os.Exit(1)
	}

	c := make(chan bool, 1)

	go func() { proxy.Listen(proxy.PxFile) }()

	logger.Logger.
		WithFields(logrus.Fields{
			"replay.scheme":                       proxy.PxFile.ReplayConfig().Scheme,
			"replay.host":                         proxy.PxFile.ReplayConfig().Host,
			"replay.port":                         proxy.PxFile.ReplayConfig().Port,
			"replay.pathRewriteSettings.strategy": proxy.PxFile.ReplayConfig().PathRewriteSettings.Strategy,
			"server.port":                         proxy.PxFile.ServerConfig().Port,
		}).Info("Proxy configuration ⚡️")

	<-c
}
