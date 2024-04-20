package http

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/collector"
	"github.com/henrywhitaker3/prompage/internal/config"
)

var (
	ErrNotFound = errors.New("service not found")
)

type Result struct {
	Query  config.Query
	Status bool
}

type ResultCache struct {
	mu        *sync.Mutex
	collector *collector.Collector
	interval  time.Duration
	time      time.Time
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

func (c *ResultCache) Get() ([]collector.Result, time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.results, c.time
}

func (c *ResultCache) GetService(name string) (collector.Result, time.Time, error) {
	results, t := c.Get()
	for _, r := range results {
		if r.Service.Name == name {
			return r, t, nil
		}
	}
	return collector.Result{}, t, ErrNotFound
}

func (c *ResultCache) Work(ctx context.Context) {
	c.mu.Lock()
	c.results = c.collector.Collect(ctx)
	c.time = time.Now()
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
			c.time = time.Now()
			c.mu.Unlock()
		}
	}
}
