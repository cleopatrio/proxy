package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cleopatrio/proxy/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (xy *Server) ReplayRequest(snapshot fiber.Ctx, proxyfile Proxyfile, path ProxyPath) error {
	if !path.EnableReplay || !xy.Proxyfile.ReplayEnabled() {
		return nil
	}

	reqTime := time.Now()
	duration := time.Duration(time.Since(reqTime))

	headers := map[string][]string{}
	for k, v := range snapshot.GetReqHeaders() {
		headers[k] = strings.Split(v, ",")
	}

	headers["Content-Type"] = []string{"application/json"}

	for _, h := range xy.Proxyfile.ReplayConfig().SuppressedHeaders {
		delete(headers, h.Name)
	}

	host := xy.Proxyfile.ReplayConfig().Host + func() string {
		port := xy.Proxyfile.ReplayConfig().Port
		if port > 0 {
			return fmt.Sprintf(":%d", port)
		}
		return ""
	}()

	reqPath := func() string {
		switch xy.Proxyfile.ReplayConfig().PathRewriteSettings.Strategy {
		case RewritePathStrategy:
			return xy.Proxyfile.ReplayConfig().PathRewriteSettings.Path
		case SuppressPathStrategy:
			return ""
		default:
			return snapshot.Path()
		}
	}()

	requestURL, err := url.Parse(xy.Proxyfile.ReplayConfig().Scheme + "://" + host + reqPath)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"request.id": snapshot.GetRespHeader(PxFile.Annotations.HTTPRequestIdHeader),
			"url":        requestURL,
			"error":      err,
		}).Error("Invalid replay URL ❌")

		return err
	}

	method := func() string {
		switch xy.Proxyfile.ReplayConfig().MethodRewriteSettings.Strategy {
		case RewriteMethodStrategy:
			return xy.Proxyfile.ReplayConfig().MethodRewriteSettings.Method
		default:
			return snapshot.Method()
		}
	}()

	data, _ := json.Marshal(map[string]any{
		"body":      snapshot.Body(),
		"path":      snapshot.Path(),
		"method":    snapshot.Method(),
		"headers":   snapshot.GetReqHeaders(),
		"remote_ip": snapshot.Context().RemoteIP(),
	})

	res, err := http.DefaultClient.Do(&http.Request{
		Method: method,
		Header: headers,
		URL:    requestURL,
		Body:   &RequestBody{Data: data},
	})

	status := -1
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"request.id": snapshot.GetRespHeader(PxFile.Annotations.HTTPRequestIdHeader),
			"url":        requestURL,
			"method":     method,
			"error":      err,
		}).Error("HTTP replay failed ❌")

		return nil
	}

	status = res.StatusCode
	defer res.Body.Close()

	fmt.Println(string(data))

	logger.Logger.WithFields(logrus.Fields{
		"request.id": snapshot.GetRespHeader(PxFile.Annotations.HTTPRequestIdHeader),
		"duration":   duration.Nanoseconds(),
		"url":        requestURL.String(),
		"method":     method,
		"status":     status,
		"error":      err,
	}).Info("Replayed HTTP request ⏪")

	return nil
}
