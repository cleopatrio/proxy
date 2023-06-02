package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cleopatrio/proxy/logger"
	"github.com/sirupsen/logrus"
)

type PathType string

const (
	ExactPathType  PathType = "Exact"
	PrefixPathType PathType = "Prefix"
)

type ProxyPath struct {
	Path     string   `yaml:"path" example:"/files"`
	PathType PathType `yaml:"pathType" example:"Exact"`
	Port     int      `yaml:"portNumber" example:"3001"`
	TLS      bool     `yaml:"tls"`
}

func (p *ProxyPath) DownstreamURL(requestHost, requestPath string) string {
	// Overwrite host scheme based on TLS configuration.
	if p.TLS {
		requestHost = strings.ReplaceAll(requestHost, "http://", "https://")
		if !strings.HasPrefix(requestHost, "https://") {
			requestHost = "https://" + requestHost
		}
	} else {
		requestHost = strings.ReplaceAll(requestHost, "https://", "http://")
		if !strings.HasPrefix(requestHost, "http://") {
			requestHost = "http://" + requestHost
		}
	}

	// Remove port. Will use `path.Port` if it is present.
	re := regexp.MustCompile(`:\d+`)
	requestHost = re.ReplaceAllLiteralString(requestHost, "")

	switch p.PathType {
	case PrefixPathType:
		if !strings.HasPrefix(requestPath, p.Path) {
			logger.Logger.
				WithFields(logrus.Fields{"host": requestHost, "path": requestPath}).
				Warn("Mismatched route")

			return ""
		}

		return func() string {
			if p.Port != 0 {
				// http[s?]://example.com[:PORT]/path.*
				return fmt.Sprintf(`%s:%d%s`, requestHost, p.Port, requestPath)
			}

			// http[s?]://example.com/path.*
			return requestHost + requestPath
		}()
	case ExactPathType:
		if requestPath != p.Path {
			logger.Logger.
				WithFields(logrus.Fields{"host": requestHost, "path": requestPath}).
				Warn("Mismatched route")

			return ""
		}

		// example.com/path
		return func() string {
			if p.Port != 0 {
				// http[s?]://example.com[:PORT]/path
				return fmt.Sprintf(`%s:%d%s`, requestHost, p.Port, p.Path)
			}

			// http[s?]://example.com/path
			return requestHost + p.Path
		}()
	}

	logger.Logger.
		WithFields(logrus.Fields{"host": requestHost, "path": requestPath}).
		Warn("Mismatched route")

	return ""
}
