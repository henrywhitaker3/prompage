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
	results := []Result{}

	for _, svc := range c.svcs {
		res := Result{
			Service: svc,
		}
		log.Printf("collecting metrics for %s\n", svc.Name)

		status, err := c.q.Status(ctx, svc.Query)
		if err != nil {
			log.Printf("ERROR - Failed to scrape status metric for %s query %s: %s", svc.Name, svc.Query.Name, err)
			res.Success = false
			res.Status = false
			results = append(results, res)
			continue
		}

		res.Success = true
		res.Status = status
		results = append(results, res)
	}

	return results
}
