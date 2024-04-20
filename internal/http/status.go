package http

import (
	"net/http"
	"time"

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

type group struct {
	Name    string
	Results []collector.Result
}

type statusData struct {
	views.Builder

	Config        config.Config
	Results       []group
	Age           time.Duration
	Outage        string
	BannerClasses string
	Version       string
	Refresh       int
}

func NewStatusPageHandler(app *app.App, cache *ResultCache) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, t := cache.Get()
		age := time.Since(t)
		op := operational(res)

		data := statusData{
			Config:        *app.Config,
			Results:       groupResults(res),
			Age:           age.Round(time.Second),
			Outage:        op,
			BannerClasses: bannerClasses(op),
			Version:       app.Version,
			Refresh:       int(app.Config.UI.RefreshInterval / time.Millisecond),
		}

		out, err := views.Build(views.STATUS, data)
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, out)
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

func groupResults(res []collector.Result) []group {
	grouped := []group{}
	// Add the default group at index 0
	grouped = append(grouped, group{Name: "default"})

	for _, r := range res {
		added := false
		for i, g := range grouped {
			if g.Name == r.Service.Group {
				g.Results = append(g.Results, r)
				grouped[i] = g
				added = true
			}
		}
		if !added {
			grouped = append(grouped, group{
				Name:    r.Service.Group,
				Results: []collector.Result{r},
			})
		}
	}
	return grouped
}
