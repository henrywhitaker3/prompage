package http

import (
	"context"
	"fmt"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type Http struct {
	e   *echo.Echo
	app *app.App
}

func NewHttp(app *app.App) *Http {
	e := echo.New()

	e.Use(echoprometheus.NewMiddleware("prompage"))

	return &Http{
		e:   e,
		app: app,
	}
}

func (h *Http) Serve() error {
	return h.e.Start(fmt.Sprintf(":%d", h.app.Config.Port))
}

func (h *Http) Stop(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}
