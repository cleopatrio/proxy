package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cleopatrio/proxy/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
)

type Host struct{ Fiber *fiber.App }

// Server - HTTP Server
type Server struct {
	App       *fiber.App
	Hosts     map[string]*Host
	Proxyfile Proxyfile
}

func (xy *Server) registerRule(rule ProxyEndpointRule) {
	app := fiber.New()

	if xy.Hosts == nil {
		xy.Hosts = map[string]*Host{}
	}

	xy.Hosts[rule.Host] = &Host{app}

	for _, path := range rule.Paths {
		/*
			Host: example.com
			Exact  -> /echo  	 -> http://example.com/echo
			Prefix -> /static/.* -> http://example.com/static/{proxy+}
		*/
		routerPath := func() string {
			switch path.PathType {
			case PrefixPathType:
				return fmt.Sprintf("%s*", path.Path)
			default:
				return path.Path
			}
		}()

		app.All(routerPath, func(c *fiber.Ctx) error {
			defer func() { go xy.ReplayRequest(*c, xy.Proxyfile, path) }()

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

func (xy *Server) getHostname(hostname string) *Host { return xy.Hosts[normalizedHostname(hostname)] }

// Listen - starts listening for HTTP requests.
func Listen(proxyfile Proxyfile) {
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

	server.Use(pprof.New(pprof.Config{Prefix: "/metrics"}))

	server.Get("/metrics", monitor.New(monitor.Config{
		Title:   "Proxy",
		FontURL: "https://fonts.googleapis.com/css2?family=REM:wght@300;400;700&display=swap",
	}))

	server.Use(RequestLoggerMiddleware)

	// ===================================
	// NOTE: Proxy routes as per Proxyfile
	// ===================================
	proxy := Server{
		App:       server,
		Hosts:     map[string]*Host{},
		Proxyfile: proxyfile,
	}

	for _, rule := range proxyfile.Rules() {
		proxy.registerRule(rule)
	}

	// TODO: Handle rate limiting
	// TODO: Handle caching

	proxy.App.Use(func(c *fiber.Ctx) error {
		if host := proxy.getHostname(c.Hostname()); host != nil {
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

	proxy.App.Listen(fmt.Sprintf(":%d", proxyfile.ServerPort()))
}

func normalizedHostname(hostname string) string {
	if components := strings.Split(hostname, ":"); len(components) > 1 {
		return components[0]
	}

	return hostname
}
