package querier

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/expr-lang/expr"
	"github.com/henrywhitaker3/prompage/internal/config"
	prometheus "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var (
	ErrTypeNotImplemented = errors.New("query result type not implemented yet")
)

type Prometheus struct {
	client v1.API
}

func NewPrometheus(conf config.Datasource) (*Prometheus, error) {
	client, err := prometheus.NewClient(prometheus.Config{
		Address: conf.Url,
		Client:  http.DefaultClient,
	})
	if err != nil {
		return nil, err
	}

	api := v1.NewAPI(client)

	return &Prometheus{
		client: api,
	}, nil
}

func (q *Prometheus) Uptime(ctx context.Context, query config.Query) (float32, []Item, error) {
	val, _, err := q.client.QueryRange(ctx, query.Query, v1.Range{
		Start: time.Now().Add(-query.Range),
		End:   time.Now(),
		Step:  query.Step,
	})
	if err != nil {
		return 0, nil, err
	}

	switch r := val.(type) {
	case model.Matrix:
		if r.Len() < 1 {
			return 0, nil, errors.New("no results for query")
		}

		passing := 0
		total := 0
		series := []Item{}
		for _, val := range r[0].Values {
			value := float64(0)
			if query.BoolValue {
				res, err := q.vector(val.Value, query)
				if err != nil {
					return 0, nil, err
				}
				if res {
					passing++
					value = 1
				}
			} else {
				f, err := q.asFloat(val.Value)
				if err != nil {
					return 0, nil, err
				}
				env := map[string]any{"result": f}
				exp, err := expr.Compile(query.Expression, expr.Env(env), expr.AsFloat64())
				if err != nil {
					return 0, nil, err
				}
				out, err := expr.Run(exp, env)
				if err != nil {
					return 0, nil, err
				}
				value = out.(float64)
			}
			total++

			series = append(series, Item{Time: val.Timestamp.Time(), Value: value})
		}

		return (float32(passing) / float32(total)) * 100, series, nil
	}

	return 100, nil, ErrTypeNotImplemented
}

func (q *Prometheus) Status(ctx context.Context, query config.Query) (bool, error) {
	val, _, err := q.client.Query(ctx, query.Query, time.Now())
	if err != nil {
		return false, err
	}

	switch r := val.(type) {
	case model.Vector:
		if r.Len() < 1 {
			return false, errors.New("no results for query")
		}
		return q.vector(r[0].Value, query)
	}

	return false, ErrTypeNotImplemented
}

func (q *Prometheus) vector(v model.SampleValue, query config.Query) (bool, error) {
	env := map[string]any{
		"result": 0,
	}

	exp, err := expr.Compile(query.Expression, expr.Env(env), expr.AsBool())
	if err != nil {
		return false, fmt.Errorf("failed to compile expr: %v", err)
	}
	val, err := q.asFloat(v)
	if err != nil {
		return false, err
	}

	env["result"] = val
	out, err := expr.Run(exp, env)
	if err != nil {
		return false, err
	}

	return out.(bool), nil
}

func (q *Prometheus) asFloat(v model.SampleValue) (float64, error) {
	val, err := strconv.ParseFloat(v.String(), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float: %v", err)
	}
	return val, nil
}
