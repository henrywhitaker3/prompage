package views

import (
	"testing"
	"time"

	"github.com/henrywhitaker3/prompage/internal/collector"
)

func TestCondensePreservesLatestSample(t *testing.T) {
	start := time.Date(2026, time.November, 7, 0, 0, 0, 0, time.UTC)
	data := make([]collector.SeriesItem, 288)
	for i := range data {
		data[i] = collector.SeriesItem{
			Time:  start.Add(time.Duration(i) * 5 * time.Minute),
			Value: float64(i),
		}
	}

	condensed := condense(collector.Series{Data: data}, 200)

	if len(condensed.Data) != 200 {
		t.Fatalf("condense() returned %d points, want 200", len(condensed.Data))
	}
	if got, want := condensed.Data[len(condensed.Data)-1].Time, data[len(data)-1].Time; !got.Equal(
		want,
	) {
		t.Fatalf("condense() ended at %v, want %v", got, want)
	}
	if got, want := condensed.Data[len(condensed.Data)-1].Value, 286.5; got != want {
		t.Fatalf("condense() returned final value %v, want %v", got, want)
	}
}

func TestCondenseLeavesSmallSeriesUnchanged(t *testing.T) {
	series := collector.Series{Data: []collector.SeriesItem{{Value: 1}, {Value: 2}}}

	condensed := condense(series, 2)

	if len(condensed.Data) != len(series.Data) {
		t.Fatalf("condense() returned %d points, want %d", len(condensed.Data), len(series.Data))
	}
	if condensed.Data[0] != series.Data[0] || condensed.Data[1] != series.Data[1] {
		t.Fatal("condense() changed a series that does not need condensing")
	}
}
