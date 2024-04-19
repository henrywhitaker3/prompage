package app

import (
	"github.com/henrywhitaker3/prompage/internal/collector"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/querier"
)

type App struct {
	Version   string
	Config    *config.Config
	Querier   *querier.Querier
	Collector *collector.Collector
}

func NewApp(conf *config.Config, q *querier.Querier) *App {
	app := &App{
		Config:  conf,
		Querier: q,
	}

	app.Collector = collector.NewCollector(app.Querier, conf.Services)

	return app
}
