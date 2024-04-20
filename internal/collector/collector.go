package collector

import (
	"context"
	"log"
	"time"

	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/metrics"
	"github.com/henrywhitaker3/prompage/internal/querier"
)

type SeriesItem struct {
	Time  time.Time
	Value float64
}

type Series struct {
	Query config.Query
	Data  []SeriesItem
}

type Result struct {
	// The service the result corresponds to
	Service config.Service
	// Whether the collection was successful or not
	Success bool
	// The boolean result of the main service query
	Status bool
	// The percentage uptime for the specified duration
	Uptime float32
	// The series of values for the range query
	Series Series
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

		status := float64(0)
		if res.Status {
			status = 1
		}
		metrics.ServiceStatus.WithLabelValues(res.Service.Name, res.Service.Group).Set(status)
		metrics.ServiceUptime.WithLabelValues(res.Service.Name, res.Service.Group).Set(float64(res.Uptime))

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
	uptime, series, err := c.q.Uptime(ctx, svc.Query)
	if err != nil {
		log.Printf("ERROR - Failed to scrape uptime metric for %s query %s: %s", svc.Name, svc.Query.Name, err)
		res.Success = false
	}

	res.Series = c.mapQuerierSeries(svc.Query, series)
	res.Status = status
	res.Uptime = uptime
	ch <- res
}

func (c *Collector) mapQuerierSeries(q config.Query, s []querier.Item) Series {
	out := Series{
		Query: q,
	}

	for _, i := range s {
		out.Data = append(out.Data, SeriesItem{
			Time:  i.Time,
			Value: i.Value,
		})
	}
	return out
}
