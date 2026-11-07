package health

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Health struct {
	e    *echo.Echo
	port int

	ready   bool
	healthy bool
}

// TODO: add an actual health/readiness check so these aren't meaningless

func NewHealth(port int) *Health {
	e := echo.New()
	h := &Health{
		e:       e,
		healthy: true,
		ready:   true,
		port:    port,
	}

	h.e.GET("/healthz", h.healthHandler())
	h.e.GET("/readyz", h.readyHandler())

	return h
}

func (h *Health) Start() error {
	return h.e.Start(fmt.Sprintf(":%d", h.port))
}

func (h *Health) healthHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		if h.healthy {
			return c.String(http.StatusOK, "HEALTHY")
		}
		return c.String(http.StatusServiceUnavailable, "UNHEALTHY")
	}
}

func (h *Health) readyHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		if h.ready {
			return c.String(http.StatusOK, "READY")
		}
		return c.String(http.StatusServiceUnavailable, "NOT READY")
	}
}
