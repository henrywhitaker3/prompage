package http

import (
	"errors"
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

type getServiceData struct {
	views.Builder

	Config  config.Config
	Age     time.Duration
	Version string
	Result  collector.Result
	Graph   template.HTML
}

func NewGetServiceHandler(app *app.App, cache *ResultCache) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		svc, age, err := cache.GetService(name)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("ERROR - could not fins service: %s\n", name)
			return err
		}
		graph, err := views.GenerateLineChart(svc.Series, app.Config.UI.Graphs.Points)
		if err != nil {
			log.Printf("ERROR - could not generate graph: %s\n", err)
		}
		data := getServiceData{
			Age:    time.Since(age).Round(time.Second),
			Result: svc,
			Graph:  template.HTML(graph),
		}

		out, err := views.Build(views.SERVICE, data)
		if err != nil {
			log.Printf("ERROR - could not render template: %s\n", err)
			return err
		}

		return c.HTML(http.StatusOK, out)
	}
}
