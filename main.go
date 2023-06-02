package main

import (
	"os"

	"github.com/cleopatrio/proxy/config"
	"github.com/cleopatrio/proxy/logger"
	"github.com/cleopatrio/proxy/proxy"
)

func main() {
	var proxyfile config.Proxyfile

	file, err := os.ReadFile("Proxyfile")
	if err != nil {
		logger.Logger.Warn("Unable to load Proxyfile. Enabled default configuration.")
	} else {
		proxyfile, err = config.LoadProxyConfiguration(file)
		if err != nil {
			logger.Logger.Fatal("Invalid Proxyfile")
		}
	}

	if proxyfile.Port() == 0 {
		proxyfile.UsePort(config.ProxyConfig.DefaultHTTPPort)

		logger.Logger.
			WithField("port", config.ProxyConfig.DefaultHTTPPort).
			Warn("Using default server port ⚠️")
	}

	c := make(chan bool, 1)

	go func() {
		proxy.Listen(proxyfile)
	}()

	logger.Logger.
		WithField("port", proxyfile.Port()).
		Info("Proxy server is running ⚡️")

	<-c
}
