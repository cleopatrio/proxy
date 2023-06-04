package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cleopatrio/proxy/core"
	"github.com/cleopatrio/proxy/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (xy *Server) ReplayRequest(snapshot *fiber.Ctx) error {
	reqTime := time.Now()
	duration := time.Duration(time.Now().Sub(reqTime))

	headers := map[string][]string{}
	for k, v := range snapshot.GetReqHeaders() {
		headers[k] = strings.Split(v, ",")
	}

	for _, h := range xy.Proxyfile.HTTPReplaySettings().SuppressHeaders {
		if _, ok := headers[h.Name]; ok {
			delete(headers, h.Name)
		}
	}

	host := xy.Proxyfile.HTTPReplaySettings().Host + func() string {
		port := xy.Proxyfile.HTTPReplaySettings().Port
		if port > 0 {
			return fmt.Sprintf(":%d", port)
		}
		return ""
	}()

	requestURL, err := url.Parse(host + snapshot.Path())
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"request.id": snapshot.GetRespHeader(core.ProxyConfig.HTTPRequestIdHeader),
			"url":        requestURL,
			"error":      err,
		}).Error("HTTP replay failed ❌")

		return err
	}

	// TODO: Perform replay...

	_ = http.Request{
		Method: snapshot.Method(),
		Header: headers,
		URL:    requestURL,
		Body:   &core.RequestBody{Data: snapshot.Body()},
	}

	logger.Logger.WithFields(logrus.Fields{
		"request.id": snapshot.GetRespHeader(core.ProxyConfig.HTTPRequestIdHeader),
		"duration":   duration.Nanoseconds(),
		"url":        requestURL.String(),
	}).Info("Replayed HTTP request ⏪")

	return err
}
