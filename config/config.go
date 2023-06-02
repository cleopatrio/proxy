package config

import (
	"sync"
)

type config struct {
	EnableRateLimiting   bool
	EnableStackTrace     bool
	EnableReplayRequests bool
	DefaultHTTPPort      int

	RequestIdHeader string
}

var once sync.Once
var ProxyConfig config

func init() {
	once.Do(func() {
		ProxyConfig = config{
			RequestIdHeader: "X-Request-Id",
			DefaultHTTPPort: 8080,
		}
	})
}
