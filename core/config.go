package core

import (
	"sync"
)

type config struct {
	EnableRateLimiting   bool
	EnableStackTrace     bool
	EnableReplayRequests bool
	DefaultHTTPPort      int

	HTTPRequestIdHeader string
}

var once sync.Once
var ProxyConfig config

func init() {
	once.Do(func() {
		ProxyConfig = config{
			HTTPRequestIdHeader: "X-Request-Id",
			DefaultHTTPPort:     8080,
		}
	})
}
