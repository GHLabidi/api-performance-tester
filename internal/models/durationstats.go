package models

import (
	"sort"

	"gonum.org/v1/gonum/stat"
)

type DurationStats struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
	Avg int64 `json:"avg"`
	Std int64 `json:"std"`
	P25 int64 `json:"p25"`
	P50 int64 `json:"p50"`
	P75 int64 `json:"p75"`
	P90 int64 `json:"p90"`
	P95 int64 `json:"p95"`
	P99 int64 `json:"p99"`
}

func CalculateStats(durations interface{}) (d DurationStats) {

	// check the type of durations
	switch durations.(type) {
	case []int64:
		// sort []int64
		sort.Slice(durations.([]int64), func(i, j int) bool { return durations.([]int64)[i] < durations.([]int64)[j] })
		// convert []int64 to []float64 (gonum/stat only works with []float64)
		var data []float64
		for _, v := range durations.([]int64) {
			data = append(data, float64(v))
		}

		d = DurationStats{
			Min: durations.([]int64)[0],
			Max: durations.([]int64)[len(durations.([]int64))-1],
			Avg: int64(stat.Mean(data, nil)),
			Std: int64(stat.StdDev(data, nil)),
			P25: int64(stat.Quantile(0.25, stat.Empirical, data, nil)),
			P50: int64(stat.Quantile(0.50, stat.Empirical, data, nil)),
			P75: int64(stat.Quantile(0.75, stat.Empirical, data, nil)),
			P90: int64(stat.Quantile(0.90, stat.Empirical, data, nil)),
			P95: int64(stat.Quantile(0.95, stat.Empirical, data, nil)),
			P99: int64(stat.Quantile(0.99, stat.Empirical, data, nil)),
		}
	case []float64:
		// sort []float64
		sort.Float64s(durations.([]float64))

		d = DurationStats{
			Min: int64(durations.([]float64)[0]),
			Max: int64(durations.([]float64)[len(durations.([]float64))-1]),
			Avg: int64(stat.Mean(durations.([]float64), nil)),
			Std: int64(stat.StdDev(durations.([]float64), nil)),
			P25: int64(stat.Quantile(0.25, stat.Empirical, durations.([]float64), nil)),
			P50: int64(stat.Quantile(0.50, stat.Empirical, durations.([]float64), nil)),
			P75: int64(stat.Quantile(0.75, stat.Empirical, durations.([]float64), nil)),
			P90: int64(stat.Quantile(0.90, stat.Empirical, durations.([]float64), nil)),
			P95: int64(stat.Quantile(0.95, stat.Empirical, durations.([]float64), nil)),
			P99: int64(stat.Quantile(0.99, stat.Empirical, durations.([]float64), nil)),
		}
	default:
		panic("unknown type")
	}
	return d
}

func CalculateStatsFromColumn(stats []StatData, column string) (d DurationStats) {

	var durations []int64
	for _, stat := range stats {
		switch column {
		case "query_duration":
			durations = append(durations, stat.QueryDuration)
		case "request_duration":
			durations = append(durations, stat.RequestDuration)
		}
	}
	return CalculateStats(durations)
}
