package http

import (
	"sync"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/querier"
)

type Result struct {
	Query  config.Query
	Status bool
}

type ResultCache struct {
	mu      *sync.Mutex
	querier *querier.Querier
	results map[string][]Result
}

func NewResultCache(app *app.App) *ResultCache {
	return &ResultCache{
		mu:      &sync.Mutex{},
		querier: app.Querier,
		results: map[string][]Result{},
	}
}
