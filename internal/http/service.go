package http

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/henrywhitaker3/flow"
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
	Extras  map[string]template.HTML
}

func NewGetServiceHandler(app *app.App, cache *ResultCache) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		svc, age, err := cache.GetService(name)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("ERROR - could not find service: %s\n", name)
			return err
		}

		group := &flow.ResultGroup{}

		graph := flow.Eventually(c.Request().Context(), func(ctx context.Context) (string, error) {
			return views.GenerateLineChart(svc.Series, app.Config.UI.Graphs.Points)
		})
		group.Add(graph)

		extraGraphs := map[string]*flow.Result[string]{}

		for _, extra := range svc.Extras {
			res := flow.Eventually(c.Request().Context(), func(ctx context.Context) (string, error) {
				return views.GenerateLineChart(extra, app.Config.UI.Graphs.Points)
			})
			group.Add(res)
			extraGraphs[extra.Query.Name] = res
		}

		// Wait for all the graphs to generate
		group.Wait()

		if graph.Err() != nil {
			log.Printf("ERROR - could not generate graph: %s\n", err)
		}

		data := getServiceData{
			Age:    time.Since(age).Round(time.Second),
			Result: svc,
			Graph:  template.HTML(graph.Out()),
			Extras: map[string]template.HTML{},
		}

		for name, res := range extraGraphs {
			if res.Err() != nil {
				log.Printf("ERROR - could not generate graph: %s\n", res.Err())
				continue
			}
			data.Extras[name] = template.HTML(res.Out())
		}

		out, err := views.Build(views.SERVICE, data)
		if err != nil {
			log.Printf("ERROR - could not render template: %s\n", err)
			return err
		}

		return c.HTML(http.StatusOK, out)
	}
}
