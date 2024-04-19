package collector

import (
	"context"
	"log"

	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/querier"
)

type Result struct {
	// The service the result corresponds to
	Service config.Service
	// Whether the collection was successful or not
	Success bool
	// The boolean result of the main service query
	Status bool
	// The percentage uptime for the specified duration
	Uptime float32
}

type Collector struct {
	q    *querier.Querier
	svcs []config.Service
}

func NewCollector(q *querier.Querier, svcs []config.Service) *Collector {
	return &Collector{
		q:    q,
		svcs: svcs,
	}
}

func (c *Collector) Collect(ctx context.Context) []Result {
	order := map[string]int{}
	results := make(chan Result, len(c.svcs))
	defer close(results)
	for i, svc := range c.svcs {
		order[svc.Name] = i
		go c.collectService(ctx, svc, results)
	}

	out := make([]Result, len(c.svcs))
	processed := 0
	for res := range results {
		out[order[res.Service.Name]] = res
		processed++
		if processed == len(c.svcs) {
			break
		}
	}

	return out
}

func (c *Collector) collectService(ctx context.Context, svc config.Service, ch chan<- Result) {
	res := Result{
		Service: svc,
		Status:  false,
		Success: true,
		Uptime:  0,
	}
	log.Printf("collecting metrics for %s\n", svc.Name)

	status, err := c.q.Status(ctx, svc.Query)
	if err != nil {
		log.Printf("ERROR - Failed to scrape status metric for %s query %s: %s", svc.Name, svc.Query.Name, err)
		res.Success = false
	}
	uptime, err := c.q.Uptime(ctx, svc.Query)
	if err != nil {
		log.Printf("ERROR - Failed to scrape uptime metric for %s query %s: %s", svc.Name, svc.Query.Name, err)
		res.Success = false
	}

	res.Success = true
	res.Status = status
	res.Uptime = uptime
	ch <- res
}
