package app

import (
	"github.com/henrywhitaker3/prompage/internal/collector"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/querier"
)

type App struct {
	Version   string
	Config    *config.Config
	Queriers  map[string]querier.Querier
	Collector *collector.Collector
}

func NewApp(conf *config.Config, q map[string]querier.Querier) *App {
	app := &App{
		Config:   conf,
		Queriers: q,
	}

	app.Collector = collector.NewCollector(conf.Services, q)

	return app
}
