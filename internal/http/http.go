package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/resources/static"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Http struct {
	e     *echo.Echo
	app   *app.App
	cache *ResultCache
}

func NewHttp(app *app.App, cache *ResultCache) *Http {
	e := echo.New()
	e.HideBanner = true

	e.Use(echoprometheus.NewMiddleware("prompage"))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Second * 5,
	}))
	if *app.Config.HttpLogs {
		e.Use(middleware.Logger())
	}

	e.GET("/", NewStatusPageHandler(app, cache))
	e.GET("/:name", NewGetServiceHandler(app, cache))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServerFS(static.FS))))

	e.GET("/api/services", NewGetAllHandler(app, cache))
	e.GET("/api/services/:name", NewGetHandler(app, cache))

	// e.GET("/metrics", echoprometheus.NewHandler())

	return &Http{
		e:     e,
		app:   app,
		cache: cache,
	}
}

func (h *Http) Start() error {
	return h.e.Start(fmt.Sprintf(":%d", h.app.Config.Port))
}

func (h *Http) Stop(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}
