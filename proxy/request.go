package proxy

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cleopatrio/proxy/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (xy *Server) MakeHTTPRequest(c *fiber.Ctx, path ProxyPath) (*http.Response, error) {
	downstreamURL := path.RequestURL(c.Hostname(), c.Path())

	if downstreamURL == nil {
		return nil, errors.New(`invalid/unknown downstream url`)
	}

	logger.Logger.
		WithFields(logrus.Fields{"method": c.Method(), "url": downstreamURL.RequestURI(), "tls": path.TLS}).
		Info("Sending HTTP request 📡")

	headers := map[string][]string{}
	for k, v := range c.GetReqHeaders() {
		headers[k] = strings.Split(v, ",")
	}

	request := http.Request{
		Method: c.Method(),
		Header: headers,
		URL:    downstreamURL,
	}

	if len(c.Body()) > 0 {
		request.Body = &RequestBody{Data: c.Body()}
	}

	return http.DefaultClient.Do(&request)
}
