package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cleopatrio/proxy/core"
	"github.com/cleopatrio/proxy/logger"
	"github.com/cleopatrio/proxy/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
)

type Host struct{ Fiber *fiber.App }

type Server struct {
	App       *fiber.App
	Hosts     map[string]*Host
	Proxyfile core.Proxyfile
}

func Listen(proxyfile core.Proxyfile) {
	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Logger.Error(err)
			c.SendStatus(http.StatusInternalServerError)
			return nil
		},
	})

	server.Use(recover.New(recover.Config{
		EnableStackTrace: proxyfile.Annotations.StackTraceEnabled,
	}))

	server.Use(requestid.New(requestid.Config{
		Header: proxyfile.Annotations.HTTPRequestIdHeader,
	}))

	server.Use(middleware.RequestLoggerMiddleware)

	// ===================================
	// NOTE: Proxy routes as per Proxyfile
	// ===================================
	proxy := Server{
		App:       server,
		Hosts:     map[string]*Host{},
		Proxyfile: proxyfile,
	}

	for _, rule := range proxyfile.Rules() {
		proxy.RegisterRule(rule)
	}

	// TODO: Handle rate limiting
	// TODO: Handle caching

	proxy.App.Use(func(c *fiber.Ctx) error {
		if host := proxy.Get(c.Hostname()); host != nil {
			logger.Logger.
				WithFields(logrus.Fields{
					"host": c.Hostname(),
					"path": c.Path(),
				}).
				Info("Handling HTTP request 📨️")

			host.Fiber.Handler()(c.Context())

			return nil
		}

		logger.Logger.
			WithFields(logrus.Fields{
				"host": c.Hostname(),
				"path": c.Path(),
			}).
			Error("Host not found 😢")

		return c.SendStatus(fiber.StatusNotFound)
	})

	proxy.App.Listen(fmt.Sprintf(":%d", proxyfile.HTTPPort()))
}

func (xy *Server) RegisterRule(rule core.ProxyRule) {
	app := fiber.New()

	if xy.Hosts == nil {
		xy.Hosts = map[string]*Host{}
	}

	xy.Hosts[rule.Host] = &Host{app}

	for _, path := range rule.Http.Paths {
		/*
			Host: example.com
			Exact  -> /echo  	 -> http://example.com/echo
			Prefix -> /static/.* -> http://example.com/static/{proxy+}
		*/
		routerPath := func() string {
			switch path.PathType {
			case core.PrefixPathType:
				return fmt.Sprintf("%s*", path.Path)
			default:
				return path.Path
			}
		}()

		app.All(routerPath, func(c *fiber.Ctx) error {
			defer func() {
				if xy.Proxyfile.Annotations.ReplayRequestsEnabled && path.EnableReplay {
					xy.ReplayRequest(c)
				}
			}()

			response, err := xy.MakeHTTPRequest(c, path)

			if err != nil {
				return c.SendStatus(http.StatusBadGateway)
			}

			c.Status(response.StatusCode)
			return c.SendStream(response.Body)
		})

		logger.Logger.WithFields(logrus.Fields{
			"host":     rule.Host,
			"pathType": path.PathType,
			"path":     path.Path,
			"port":     path.PortNumber,
			"tls":      path.TLS,
		}).Debug("Registered route")
	}
}

func (xy *Server) MakeHTTPRequest(c *fiber.Ctx, path core.ProxyPath) (*http.Response, error) {
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
		request.Body = &core.RequestBody{Data: c.Body()}
	}

	return http.DefaultClient.Do(&request)
}

func (xy *Server) Get(hostname string) *Host { return xy.Hosts[normalizedHostname(hostname)] }

func normalizedHostname(hostname string) string {
	if components := strings.Split(hostname, ":"); len(components) > 1 {
		return components[0]
	}

	return hostname
}
