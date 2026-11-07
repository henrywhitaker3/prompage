package views

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/henrywhitaker3/prompage/internal/collector"
)

type chartItem struct {
	data  string
	until time.Time
}

type chartCache struct {
	mu     *sync.Mutex
	charts map[string]chartItem
}

func (c chartCache) get(name string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.charts[name]
	if !ok {
		return "", false
	}
	if time.Now().After(item.until) {
		delete(c.charts, name)
		return "", false
	}
	return item.data, true
}

func (c *chartCache) put(name string, data string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.charts[name] = chartItem{
		data:  data,
		until: time.Now().Add(time.Minute),
	}
}

func cacheKey(series collector.Series) string {
	return fmt.Sprintf("%s:%s:%s", series.Query.Name, series.Query.Query, series.Query.Expression)
}

var (
	cache = &chartCache{
		mu:     &sync.Mutex{},
		charts: map[string]chartItem{},
	}
)

// Generates a line chart in html for a given series
func GenerateLineChart(series collector.Series, maxPoints int) (string, error) {
	if chart, ok := cache.get(cacheKey(series)); ok {
		return chart, nil
	}

	line := charts.NewLine()

	yopts := opts.YAxis{
		Show: opts.Bool(true),
		Type: "value",
	}
	if series.Query.BoolValue {
		yopts.Max = 2
		yopts.Min = 0
	}
	if series.Query.Units != "" {
		yopts.Name = series.Query.Units
	}

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros,
		}),
		charts.WithYAxisOpts(yopts),
		charts.WithXAxisOpts(opts.XAxis{
			Show: opts.Bool(true),
			Type: "time",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)

	cd := condense(series, maxPoints)

	line.SetXAxis(nil).
		AddSeries("Metric", getYAxis(cd)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth:     opts.Bool(true),
				ShowSymbol: opts.Bool(false),
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{Opacity: opts.Float(0.2), Color: "#5c6848"}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color: "#5c6848",
			}),
		)

	cr := newChartRenderer(line, line.Validate)
	var out bytes.Buffer
	if err := cr.Render(&out); err != nil {
		return "", err
	}
	chart := out.String()
	cache.put(cacheKey(series), chart)

	return chart, nil
}

func getYAxis(series collector.Series) []opts.LineData {
	out := []opts.LineData{}
	for _, i := range series.Data {
		out = append(out, opts.LineData{
			Value: []interface{}{i.Time.Format(time.RFC3339), i.Value},
		})
	}
	return out
}

// Condenses the series based on the configured config ui.graphs.points
func condense(series collector.Series, to int) collector.Series {
	if to <= 0 || len(series.Data) <= to {
		return series
	}

	final := make([]collector.SeriesItem, to)
	for i := range final {
		start := i * len(series.Data) / to
		end := (i + 1) * len(series.Data) / to
		final[i] = average(series.Data[start:end])
	}
	series.Data = final

	return series
}

func average(items []collector.SeriesItem) collector.SeriesItem {
	sum := float64(0)
	end := items[len(items)-1].Time

	for _, item := range items {
		sum += item.Value
	}

	return collector.SeriesItem{
		Time:  end,
		Value: sum / float64(len(items)),
	}
}
