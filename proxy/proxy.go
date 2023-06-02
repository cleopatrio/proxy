package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cleopatrio/proxy/config"
	"github.com/cleopatrio/proxy/logger"
	"github.com/cleopatrio/proxy/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
)

func Listen(proxyfile config.Proxyfile) {
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
		Header: proxyfile.Annotations.RequestIdHeader,
	}))

	server.Use(middleware.RequestLoggerMiddleware)

	// ===================================
	// NOTE: Proxy routes as per Proxyfile
	// ===================================
	proxy := Proxy{
		Server: server,
		Hosts:  map[string]*Host{},
	}

	for _, rule := range proxyfile.Rules() {
		proxy.RegisterRule(rule)
	}

	// TODO: Handle request replay
	// TODO: Handle rate limiting
	// TODO: Handle caching

	proxy.Server.Use(func(c *fiber.Ctx) error {
		if host := proxy.Get(c.Hostname()); host != nil {
			logger.Logger.
				WithFields(logrus.Fields{
					"host": c.Hostname(),
					"path": c.Path(),
				}).
				Info("Handling path for host ⚡️")

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

	proxy.Server.Listen(fmt.Sprintf(":%d", proxyfile.Port()))
}

type Host struct{ Fiber *fiber.App }

type Proxy struct {
	Server *fiber.App
	Hosts  map[string]*Host
}

func (xy *Proxy) RegisterRule(rule config.ProxyRule) {
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
			case config.PrefixPathType:
				return fmt.Sprintf("%s*", path.Path)
			default:
				return path.Path
			}
		}()

		app.All(routerPath, func(c *fiber.Ctx) error {
			response, err := xy.MakeHTTPRequest(c, path)
			if err != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}

			c.Status(response.StatusCode)
			return c.SendStream(response.Body)
		})

		logger.Logger.WithFields(logrus.Fields{
			"host":     rule.Host,
			"pathType": path.PathType,
			"path":     path.Path,
			"port":     path.Port,
			"tls":      path.TLS,
		}).Debug("Registered route")
	}
}

func (xy *Proxy) MakeHTTPRequest(c *fiber.Ctx, path config.ProxyPath) (*http.Response, error) {
	downstreamURL := path.DownstreamURL(c.Hostname(), c.Path())

	logger.Logger.
		WithFields(logrus.Fields{
			"method": c.Method(),
			"url":    downstreamURL,
			"tls":    path.TLS,
		}).
		Info("Forwarding request downstream")

	switch c.Method() {
	case "GET":
		return http.Get(downstreamURL)
	case "POST":
		return http.Post(downstreamURL, c.Get("Content-Type"), c.Request().BodyStream())
	default:
		// TODO: Implement other methods
	}

	return nil, nil
}

func (xy *Proxy) Get(hostname string) *Host {
	return xy.Hosts[normalizedHostname(hostname)]
}

func normalizedHostname(hostname string) string {
	components := strings.Split(hostname, ":")
	if len(components) > 1 {
		return components[0]
	}

	return hostname
}
