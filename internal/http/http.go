package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/resources/static"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type Http struct {
	e   *echo.Echo
	app *app.App
}

func NewHttp(app *app.App) *Http {
	e := echo.New()
	e.HideBanner = true

	e.Use(echoprometheus.NewMiddleware("prompage"))

	e.GET("/", NewStatusPageHandler(app))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServerFS(static.FS))))

	return &Http{
		e:   e,
		app: app,
	}
}

func (h *Http) Start() error {
	return h.e.Start(fmt.Sprintf(":%d", h.app.Config.Port))
}

func (h *Http) Stop(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}
