package views

import (
	"bytes"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/henrywhitaker3/prompage/internal/collector"
)

// Generates a lin chart ing html for a given series
func GenerateLineChart(series collector.Series, maxPoints int) (string, error) {
	line := charts.NewLine()

	yopts := opts.YAxis{
		Show: false,
	}
	if series.Query.BoolValue {
		yopts.Max = 2
		yopts.Min = 0
	}

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros,
		}),
		charts.WithYAxisOpts(yopts),
		charts.WithXAxisOpts(opts.XAxis{Show: false}),
		charts.WithLegendOpts(opts.Legend{
			Show: false,
		}),
	)

	cd := condense(series, maxPoints)

	line.SetXAxis(getXAxis(cd)).
		AddSeries("Metric", getYAxis(cd)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: true, Color: "#5c6848"}),
			charts.WithAreaStyleOpts(opts.AreaStyle{Opacity: 0.2, Color: "#5c6848"}),
		)

	cr := newChartRenderer(line, line.Validate)
	var out bytes.Buffer
	if err := cr.Render(&out); err != nil {
		return "", err
	}
	return out.String(), nil
}

func getXAxis(series collector.Series) []string {
	out := []string{}
	for _, i := range series.Data {
		out = append(out, i.Time.String())
	}
	return out
}

func getYAxis(series collector.Series) []opts.LineData {
	out := []opts.LineData{}
	for _, i := range series.Data {
		out = append(out, opts.LineData{
			Value: i.Value,
		})
	}
	return out
}

// Condenses the series based on the configured config ui.graphs.points
func condense(series collector.Series, to int) collector.Series {
	if len(series.Data) <= to {
		return series
	}

	perBucket := len(series.Data) / to
	chunks := chunk(series.Data, perBucket, to)

	final := make([]collector.SeriesItem, to)
	for i := range chunks {
		final[i] = average(chunks[i])
	}
	series.Data = final

	return series
}

func chunk(items []collector.SeriesItem, perBucket, buckets int) [][]collector.SeriesItem {
	out := make([][]collector.SeriesItem, buckets)

	feed := make(chan collector.SeriesItem, 1)
	run := true
	go func() {
		for _, item := range items {
			if !run {
				break
			}
			feed <- item
		}
	}()

	for i := 0; i < buckets; i++ {
		for range perBucket {
			out[i] = append(out[i], <-feed)
		}
	}

	// Pull the last item out so it doesn't leak goroutines every graph load
	run = false
	<-feed

	return out
}

func average(items []collector.SeriesItem) collector.SeriesItem {
	sum := float64(0)
	start := items[0].Time
	end := items[len(items)-1].Time
	mid := start.Add(time.Second * (time.Duration(end.Unix() - start.Unix())))

	for _, item := range items {
		sum += item.Value
	}

	return collector.SeriesItem{
		Time:  mid,
		Value: sum / float64(len(items)),
	}
}
