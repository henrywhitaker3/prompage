package http

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/collector"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/resources/views"
	"github.com/labstack/echo/v4"
)

func NewStatusPageHandler(app *app.App, cache *ResultCache) echo.HandlerFunc {
	tmpl := template.Must(template.ParseFS(views.Views, "index.html"))

	return func(c echo.Context) error {
		res, t := cache.Get()
		age := time.Since(t)

		data := struct {
			Config  config.Config
			Results []collector.Result
			Age     time.Duration
		}{
			Config:  *app.Config,
			Results: res,
			Age:     age.Round(time.Second),
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			log.Printf("ERROR - could not render template: %s", err)
			return err
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}
