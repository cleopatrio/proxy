package main

import (
	"os"

	"github.com/cleopatrio/proxy/core"
	"github.com/cleopatrio/proxy/logger"
	"github.com/cleopatrio/proxy/proxy"
)

func main() {
	var proxyfile core.Proxyfile

	file, err := os.ReadFile("Proxyfile")
	if err != nil {
		logger.Logger.Warn("Unable to load Proxyfile. Enabled default configuration.")
	} else {
		if proxyfile, err = core.LoadProxyConfiguration(file); err != nil {
			logger.Logger.Fatal("Invalid Proxyfile ", err)
		}
	}

	c := make(chan bool, 1)

	go func() {
		proxy.Listen(proxyfile)
	}()

	logger.Logger.
		WithField("port", proxyfile.HTTPPort()).
		Info("Proxy server is running ⚡️")

	<-c
}
