package http

import (
	"context"
	"fmt"

	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type Http struct {
	e    *echo.Echo
	conf *config.Config
}

func NewHttp(conf *config.Config) *Http {
	e := echo.New()

	e.Use(echoprometheus.NewMiddleware("prompage"))

	return &Http{
		e:    e,
		conf: conf,
	}
}

func (h *Http) Serve() error {
	return h.e.Start(fmt.Sprintf(":%d", h.conf.Port))
}

func (h *Http) Stop(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}
