package config

import (
	"gopkg.in/yaml.v3"
)

type (
	Proxyfile struct {
		Annotations struct {
			RateLimitingEnabled   bool   `yaml:"proxy.conf/rate-limiting-enabled"`
			ReplayRequestsEnabled bool   `yaml:"proxy.conf/replay-requests-enabled"`
			RequestIdHeader       string `yaml:"proxy.conf/request-id-header"`
			StackTraceEnabled     bool   `yaml:"proxy.conf/stack-trace-enabled"`
		} `yaml:"annotations"`

		Spec ProxySpec `yaml:"spec"`
	}

	ProxySpec struct {
		Rules  []ProxyRule `yaml:"rules"`
		Server struct {
			Port int `yaml:"port"`
		} `yaml:"server"`
	}

	ProxyRule struct {
		Host string `yaml:"host"`
		Http struct {
			Paths []ProxyPath `yaml:"paths"`
		} `yaml:"http"`
	}
)

func LoadProxyConfiguration(file []byte) (proxy Proxyfile, err error) {
	proxy.Annotations.RateLimitingEnabled = ProxyConfig.EnableRateLimiting
	proxy.Annotations.ReplayRequestsEnabled = ProxyConfig.EnableReplayRequests
	proxy.Annotations.RequestIdHeader = ProxyConfig.RequestIdHeader
	proxy.Annotations.StackTraceEnabled = ProxyConfig.EnableStackTrace

	err = yaml.Unmarshal(file, &proxy)
	if proxy.Annotations.RequestIdHeader == "" {
		proxy.Annotations.RequestIdHeader = ProxyConfig.RequestIdHeader
	}

	return
}

func (pf *Proxyfile) Port() int { return pf.Spec.Server.Port }

func (pf *Proxyfile) UsePort(port int) { pf.Spec.Server.Port = port }

func (pf *Proxyfile) Rules() []ProxyRule { return pf.Spec.Rules }
