package middleware

import (
	"encoding/json"
	"time"

	"github.com/cleopatrio/proxy/config"
	"github.com/cleopatrio/proxy/logger"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

func RequestLoggerMiddleware(c *fiber.Ctx) error {
	reqStart := time.Now()

	// Capture any error returned by the handler
	err := c.Next()

	reqEnd := time.Now()

	var reqBody map[string]any
	json.Unmarshal(c.Body(), &reqBody)

	var resBody map[string]any
	json.Unmarshal(c.Response().Body(), &resBody)

	var headers map[string]string
	if rawHeaders, err := json.Marshal(c.GetReqHeaders()); err == nil {
		json.Unmarshal(rawHeaders, &headers)
		if _, ok := headers["Authorization"]; ok {
			headers["Authorization"] = "***"
		}
	}

	_ = map[string]any{
		// 1. Proxy Config
		"proxy.replay_requests_enabled": config.ProxyConfig.EnableReplayRequests,
		"proxy.rate_limiting_enabled":   config.ProxyConfig.EnableRateLimiting,
		"proxy.stack_trace_enabled":     config.ProxyConfig.EnableStackTrace,
		"proxy.request_id_header":       config.ProxyConfig.RequestIdHeader,

		// 2. Request
		"request.id":      c.GetRespHeader(config.ProxyConfig.RequestIdHeader),
		"request.headers": headers,
		"request.method":  c.Method(),
		"request.url":     c.BaseURL(),
		"request.path":    c.Path(),
		"request.params":  c.AllParams(),
		"request.body":    reqBody,

		// 3. Response
		"response.status":   c.Response().StatusCode(),
		"response.body":     resBody,
		"response.headers":  c.GetRespHeaders(),
		"response.duration": time.Duration(reqEnd.Sub(reqStart).Milliseconds()),
	}

	logger.Logger.
		WithFields(logrus.Fields{
			"port":   c.Port(),
			"ip":     c.IP(),
			"method": c.Method(),
			"status": c.Response().StatusCode(),
			"path":   c.Path(),
			"url":    c.BaseURL(),
		}).
		WithContext(c.UserContext()).
		Info("Handled incoming HTTP request")

	return err
}
