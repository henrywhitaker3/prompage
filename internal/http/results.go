package http

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/collector"
	"github.com/henrywhitaker3/prompage/internal/config"
)

type Result struct {
	Query  config.Query
	Status bool
}

type ResultCache struct {
	mu        *sync.Mutex
	collector *collector.Collector
	interval  time.Duration
	results   []collector.Result
}

func NewResultCache(app *app.App) *ResultCache {
	return &ResultCache{
		mu:        &sync.Mutex{},
		collector: app.Collector,
		interval:  app.Config.Refresh,
		results:   []collector.Result{},
	}
}

func (c *ResultCache) Get() []collector.Result {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.results
}

func (c *ResultCache) Work(ctx context.Context) {
	c.mu.Lock()
	c.results = c.collector.Collect(ctx)
	c.mu.Unlock()

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	log.Printf("collecitng metrics every %s", c.interval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			c.results = c.collector.Collect(ctx)
			c.mu.Unlock()
		}
	}
}
