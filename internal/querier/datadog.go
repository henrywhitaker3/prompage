package querier

import (
	"context"
	"errors"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/expr-lang/expr"
	"github.com/henrywhitaker3/prompage/internal/config"
)

var (
	ErrInvalidNumberOfSeries = errors.New("datadog returned an invalid number of series")
)

type Datadog struct {
	// The dd api client
	client *datadog.APIClient
	// The dd metrics client
	metrics *datadogV1.MetricsApi
	// The base ctx to make requests from
	ctx context.Context
}

func NewDatadog(cfg config.Datasource) (*Datadog, error) {
	if cfg.Type != "datadog" {
		return nil, ErrInvalidDatasourceMapping
	}

	client := datadog.NewAPIClient(datadog.NewConfiguration())
	metrics := datadogV1.NewMetricsApi(client)

	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: cfg.Extras["apiKey"],
			},
			"appKeyAuth": {
				Key: cfg.Extras["appKey"],
			},
		},
	)
	ctx = context.WithValue(
		ctx,
		datadog.ContextServerVariables,
		map[string]string{"site": cfg.Url},
	)

	return &Datadog{
		client:  client,
		metrics: metrics,
		ctx:     ctx,
	}, nil
}

func (d *Datadog) Uptime(ctx context.Context, query config.Query) (float32, []Item, error) {
	resp, _, err := d.metrics.QueryMetrics(
		d.ctx,
		time.Now().Add(-query.Range).Unix(),
		time.Now().Unix(),
		query.Query,
	)

	if err != nil {
		return 0, nil, err
	}

	if len(resp.Series) != 1 {
		return 0, nil, ErrInvalidNumberOfSeries
	}

	series := resp.Series[0]

	items := []Item{}
	passing := 0
	total := 0

	for _, point := range series.GetPointlist() {
		time := time.Unix(0, int64(*point[0]*float64(time.Millisecond)))
		raw := *point[1]
		value := float64(0)

		total++
		env := map[string]any{
			"result": raw,
		}
		if query.BoolValue {
			exp, err := expr.Compile(query.Expression, expr.Env(env), expr.AsBool())
			if err != nil {
				return 0, nil, err
			}
			out, err := expr.Run(exp, env)
			if err != nil {
				return 0, nil, err
			}
			if out.(bool) {
				value = 1
				passing++
			}
		} else {
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

		items = append(items, Item{
			Time:  time,
			Value: value,
		})
	}

	return (float32(passing) / float32(total)) * 100, items, nil
}

func (d *Datadog) Status(ctx context.Context, query config.Query) (bool, error) {
	resp, _, err := d.metrics.QueryMetrics(
		d.ctx,
		time.Now().Add(-query.Range).Unix(),
		time.Now().Unix(),
		query.Query,
	)
	if err != nil {
		return false, err
	}

	if len(resp.Series) != 1 {
		return false, ErrInvalidNumberOfSeries
	}

	series := resp.Series[0]

	env := map[string]any{"result": 0}
	exp, err := expr.Compile(query.Expression, expr.Env(env), expr.AsBool())
	if err != nil {
		return false, err
	}

	latest := series.GetPointlist()[series.GetLength()-1]

	env["result"] = *latest[1]
	out, err := expr.Run(exp, env)
	if err != nil {
		return false, err
	}

	return out.(bool), nil
}
