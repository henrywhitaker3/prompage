package querier

import (
	"context"
	"errors"
	"time"

	"github.com/henrywhitaker3/prompage/internal/config"
)

var (
	ErrInvalidDatasourceMapping = errors.New("unknown datasource/querier mapping")
)

type Querier interface {
	// Calculate the uptime of the service
	// Returns the % uptime, and the series of items of the metric
	Uptime(context.Context, config.Query) (float32, []Item, error)

	// Check whether the service is up/down
	Status(context.Context, config.Query) (bool, error)
}

type Item struct {
	Time  time.Time
	Value float64
}

func BuildQueriers(sources []config.Datasource) (map[string]Querier, error) {
	out := map[string]Querier{}

	for _, ds := range sources {
		switch ds.Type {
		case "prometheus":
			q, err := NewPrometheus(ds)
			if err != nil {
				return nil, err
			}
			out[ds.Name] = q
		case "datadog":
			return nil, errors.New("datadog not implemented yet")
		default:
			return nil, ErrInvalidDatasourceMapping
		}
	}

	return out, nil
}
