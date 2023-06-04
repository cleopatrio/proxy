package middleware

import (
	"time"

	"github.com/cleopatrio/proxy/core"
	"github.com/cleopatrio/proxy/logger"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

func RequestLoggerMiddleware(c *fiber.Ctx) error {
	reqStart := time.Now()

	// Capture any error returned by the handler
	err := c.Next()

	duration := time.Duration(time.Now().Sub(reqStart))

	logger.Logger.
		WithFields(logrus.Fields{
			"request.id": c.GetRespHeader(core.ProxyConfig.HTTPRequestIdHeader),
			"port":       c.Port(),
			"ip":         c.IP(),
			"method":     c.Method(),
			"status":     c.Response().StatusCode(),
			"path":       c.Path(),
			"url":        c.BaseURL(),
			"duration":   duration.Nanoseconds(),
		}).
		WithContext(c.UserContext()).
		Info("HTTP request finished ✅")

	return err
}
