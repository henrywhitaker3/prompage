package http

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/hako/durafmt"
	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/collector"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/resources/views"
	"github.com/labstack/echo/v4"
)

var (
	OutageNone    = "None"
	OutagePartial = "Partial"
	OutageFull    = "Full"
)

type statusData struct {
	Config        config.Config
	Results       []collector.Result
	Age           time.Duration
	Outage        string
	BannerClasses string
	Version       string
}

func (s statusData) Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func (s statusData) PrettyDuration(duration time.Duration) string {
	return durafmt.Parse(duration).String()
}

func NewStatusPageHandler(app *app.App, cache *ResultCache) echo.HandlerFunc {
	tmpl := template.Must(template.ParseFS(views.Views, "index.html"))

	return func(c echo.Context) error {
		res, t := cache.Get()
		age := time.Since(t)
		op := operational(res)

		data := statusData{
			Config:        *app.Config,
			Results:       res,
			Age:           age.Round(time.Second),
			Outage:        op,
			BannerClasses: bannerClasses(op),
			Version:       app.Version,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			log.Printf("ERROR - could not render template: %s", err)
			return err
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}

func operational(res []collector.Result) string {
	passing := 0
	for _, r := range res {
		if r.Status {
			passing++
		}
	}

	switch passing {
	case 0:
		return OutageFull
	case len(res):
		return OutageNone
	default:
		return OutagePartial
	}
}

func bannerClasses(outage string) string {
	switch outage {
	case OutageNone:
		return "bg-lime-600 text-white"
	case OutageFull:
		return "bg-red-500 text-white"
	case OutagePartial:
		fallthrough
	default:
		return "bg-orange-400"
	}
}
