package http

import (
	"errors"
	"net/http"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/collector"
	"github.com/labstack/echo/v4"
)

type HttpResult struct {
	Name   string  `json:"name"`
	Group  string  `json:"group"`
	Status bool    `json:"status"`
	Uptime float32 `json:"uptime"`
}

type GetAllResponse struct {
	Services []HttpResult `json:"services"`
}

func NewGetAllHandler(app *app.App, cache *ResultCache) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, _ := cache.Get()

		out := []HttpResult{}
		for _, r := range res {
			out = append(out, convertResult(r))
		}

		return c.JSON(http.StatusOK, &GetAllResponse{Services: out})
	}
}

func NewGetHandler(app *app.App, cache *ResultCache) echo.HandlerFunc {
	return func(c echo.Context) error {
		svc, _, err := cache.GetService(c.Param("name"))
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return c.JSON(http.StatusNotFound, struct{}{})
			}
			return err
		}

		return c.JSON(http.StatusOK, convertResult(svc))
	}
}

func convertResult(r collector.Result) HttpResult {
	return HttpResult{
		Name:   r.Service.Name,
		Group:  r.Service.Group,
		Status: r.Status,
		Uptime: r.Uptime,
	}
}
