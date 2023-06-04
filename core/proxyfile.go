package core

import (
	"github.com/cleopatrio/proxy/logger"
	"gopkg.in/yaml.v3"
)

type (
	// Proxy configuration.
	Proxyfile struct {
		Annotations struct {
			// Globally enables or disables request playback.
			// 	- If this setting is disabled, no requests will be replayed.
			// 	- If this setting is enabled, each path can enable or disable this behavior.
			ReplayRequestsEnabled bool `yaml:"proxy.conf/replay-requests-enabled"`

			// Globally enables or disables rate limiting.
			// 	- If this setting is disabled, no rate limiting will be applied.
			// 	- If this setting is enabled, each path can enable or disable this behavior.
			RateLimitingEnabled bool `yaml:"proxy.conf/rate-limiting-enabled"`

			// A header that uniquely identifies each incoming request.
			HTTPRequestIdHeader string `yaml:"proxy.conf/request-id-header"`

			// Dump the application stack trace if/when unexpected server errors occur.
			StackTraceEnabled bool `yaml:"proxy.conf/stack-trace-enabled"`
		} `yaml:"annotations"`

		Spec ProxySpec `yaml:"spec"`
	}

	ProxySpec struct {
		Rules  []ProxyRule `yaml:"rules"`
		Server struct {
			HTTP struct {
				// Server port.
				Port int `yaml:"port"`
				// Server replay settings.
				Replay Replay `yaml:"replay"`
			} `yaml:"http"`
		} `yaml:"server"`
	}

	ProxyRule struct {
		Host string `yaml:"host"`
		Http struct {
			Paths []ProxyPath `yaml:"paths"`
		} `yaml:"http"`
	}

	// Controls where and how HTTP requests are replayed
	Replay struct {
		// Replayed requests will be sent to this host.
		Host string `yaml:"host"`
		// Replayed requests will be sent to this port.
		Port int `yaml:"port"`
		// Replayed requests will not include these headers.
		SuppressHeaders []struct{ Name string } `yaml:"suppressHeaders"`
	}
)

func LoadProxyConfiguration(file []byte) (proxyfile Proxyfile, err error) {
	proxyfile.Annotations.HTTPRequestIdHeader = ProxyConfig.HTTPRequestIdHeader
	proxyfile.Annotations.RateLimitingEnabled = ProxyConfig.EnableRateLimiting
	proxyfile.Annotations.ReplayRequestsEnabled = ProxyConfig.EnableReplayRequests
	proxyfile.Annotations.StackTraceEnabled = ProxyConfig.EnableStackTrace

	err = yaml.Unmarshal(file, &proxyfile)
	if proxyfile.Annotations.HTTPRequestIdHeader == "" {
		proxyfile.Annotations.HTTPRequestIdHeader = ProxyConfig.HTTPRequestIdHeader
	}

	if proxyfile.HTTPPort() == 0 {
		proxyfile.UseHTTPPort(ProxyConfig.DefaultHTTPPort)

		logger.Logger.
			WithField("port", ProxyConfig.DefaultHTTPPort).
			Warn("Using default server port ⚠️")
	}

	return
}

func (pf *Proxyfile) HTTPReplaySettings() Replay { return pf.Spec.Server.HTTP.Replay }

func (pf *Proxyfile) HTTPPort() int { return pf.Spec.Server.HTTP.Port }

func (pf *Proxyfile) UseHTTPPort(port int) { pf.Spec.Server.HTTP.Port = port }

func (pf *Proxyfile) Rules() []ProxyRule { return pf.Spec.Rules }
