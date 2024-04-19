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

type Querier struct {
	client v1.API
}

type Result struct {
	Result any
	Time   time.Time
}

func NewQuerier(conf *config.Config) (*Querier, error) {
	client, err := prometheus.NewClient(prometheus.Config{
		Address: conf.Prometheus,
		Client:  http.DefaultClient,
	})
	if err != nil {
		return nil, err
	}

	api := v1.NewAPI(client)

	return &Querier{
		client: api,
	}, nil
}

func (q *Querier) Range(ctx context.Context, query config.Query) (*Result, error) {
	return nil, errors.New("range queries not implemented yet")
}

func (q *Querier) Status(ctx context.Context, query config.Query) (bool, error) {
	val, _, err := q.client.Query(ctx, query.Query, time.Now())
	if err != nil {
		return false, err
	}

	env := map[string]any{
		"result": 0,
	}

	exp, err := expr.Compile(query.Expression, expr.Env(env), expr.AsBool())
	if err != nil {
		return false, fmt.Errorf("failed to compile expr: %v", err)
	}

	switch r := val.(type) {
	case model.Vector:
		if r.Len() < 1 {
			return false, errors.New("no results for query")
		}
		// if r.Len() != 1 {
		// 	return nil, errors.New("unexpected result length")
		// }

		val, err := strconv.ParseFloat(r[0].Value.String(), 64)
		if err != nil {
			return false, fmt.Errorf("failed to parse result: %v", err)
		}

		env["result"] = val

		out, err := expr.Run(exp, env)
		if err != nil {
			return false, err
		}

		return out.(bool), nil
	}

	return false, nil
}
