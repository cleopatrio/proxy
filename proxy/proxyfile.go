package proxy

import (
	"sync"
)

const (
	// Request path
	PreservePathStrategy PathRewriteStrategy = "preserve"
	RewritePathStrategy  PathRewriteStrategy = "rewrite"
	SuppressPathStrategy PathRewriteStrategy = "suppress"

	// Request method
	PreserveMethodStrategy MethodRewriteStrategy = "preserve"
	RewriteMethodStrategy  MethodRewriteStrategy = "rewrite"

	// Server defaults
	DefaultHTTPPort      int    = 8080
	EnableRateLimiting   bool   = false
	EnableStackTrace     bool   = false
	EnableReplayRequests bool   = false
	HTTPRequestIdHeader  string = "X-Request-Id"
)

var (
	once   sync.Once
	PxFile Proxyfile
)

// PathRewriteStrategy - Controls whether the original request path should be preserved
type PathRewriteStrategy string

// MethodRewriteStrategy - Controls whether the original request method should be preserved
type MethodRewriteStrategy string

// Proxifyle - Proxy configuration.
type Proxyfile struct {
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

// ProxySpec - Proxy server and rules configuration.
type ProxySpec struct {
	Rules  []ProxyEndpointRule `yaml:"rules"`
	Server ProxyServer         `yaml:"server"`
}

// ProxyServer - Proxy server configuration.
type ProxyServer struct {
	// Server port.
	Port int `yaml:"port"`
	// Server replay settings.
	Replay ProxyReplay `yaml:"replay"`
}

// ProxyEndpointRule - Endpoint route configuration.
type ProxyEndpointRule struct {
	Host  string      `yaml:"host"`
	Paths []ProxyPath `yaml:"paths"`
}

// ProxyReplay - Controls where and how HTTP requests are replayed
type ProxyReplay struct {
	// Replayed requests will be sent using this protocol [http/https]
	Scheme string `yaml:"scheme"`

	// Replayed requests will be sent to this host.
	Host string `yaml:"host"`

	// Replayed requests will be sent to this port.
	Port int `yaml:"port"`

	// Replayed requests will not include these headers.
	SuppressedHeaders []struct{ Name string } `yaml:"suppressedHeaders"`

	MethodRewriteSettings struct {
		Strategy MethodRewriteStrategy
		Method   string
	} `yaml:"methodRewriteSettings"`

	PathRewriteSettings struct {
		Strategy PathRewriteStrategy
		Path     string
	} `yaml:"pathRewriteSettings"`
}

func (pf *Proxyfile) ReplayConfig() ProxyReplay { return pf.Spec.Server.Replay }

func (pf *Proxyfile) ServerConfig() ProxyServer { return pf.Spec.Server }

func (pf *Proxyfile) ServerPort() int { return pf.Spec.Server.Port }

func (pf *Proxyfile) UseServerPort(port int) { pf.Spec.Server.Port = port }

func (pf *Proxyfile) Rules() []ProxyEndpointRule { return pf.Spec.Rules }

func (pf *Proxyfile) ReplayEnabled() bool { return pf.Annotations.ReplayRequestsEnabled }

func init() {
	once.Do(func() {
		PxFile.Spec.Server.Port = DefaultHTTPPort
		PxFile.Annotations.HTTPRequestIdHeader = HTTPRequestIdHeader
		PxFile.Annotations.ReplayRequestsEnabled = EnableReplayRequests
		PxFile.Annotations.RateLimitingEnabled = EnableRateLimiting
		PxFile.Annotations.StackTraceEnabled = EnableStackTrace

		PxFile.Spec.Server.Replay.MethodRewriteSettings.Strategy = PreserveMethodStrategy
		PxFile.Spec.Server.Replay.PathRewriteSettings.Strategy = PreservePathStrategy
		PxFile.Spec.Server.Replay.Scheme = "http"
	})
}
