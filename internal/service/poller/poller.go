package poller

import (
	"math/rand/v2"
	"runtime"
	"strconv"

	"github.com/mixailo/go-training-metrics/internal/service/metrics"
)

func PollMetrics() metrics.Report {
	report := metrics.NewReport()
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	report.AddUnConverted(metrics.TypeGauge, "Alloc", strconv.FormatUint(memStats.Alloc, 10))
	report.AddUnConverted(metrics.TypeGauge, "Alloc", strconv.FormatUint(memStats.Alloc, 10))
	report.AddUnConverted(metrics.TypeGauge, "BuckHashSys", strconv.FormatUint(memStats.BuckHashSys, 10))
	report.AddUnConverted(metrics.TypeGauge, "Frees", strconv.FormatUint(memStats.Frees, 10))
	report.AddUnConverted(metrics.TypeGauge, "GCCPUFraction", strconv.FormatFloat(memStats.GCCPUFraction, 'f', -1, 64))
	report.AddUnConverted(metrics.TypeGauge, "GCSys", strconv.FormatUint(memStats.GCSys, 10))
	report.AddUnConverted(metrics.TypeGauge, "HeapAlloc", strconv.FormatUint(memStats.HeapAlloc, 10))
	report.AddUnConverted(metrics.TypeGauge, "HeapIdle", strconv.FormatUint(memStats.HeapIdle, 10))
	report.AddUnConverted(metrics.TypeGauge, "HeapInuse", strconv.FormatUint(memStats.HeapInuse, 10))
	report.AddUnConverted(metrics.TypeGauge, "HeapObjects", strconv.FormatUint(memStats.HeapObjects, 10))
	report.AddUnConverted(metrics.TypeGauge, "HeapReleased", strconv.FormatUint(memStats.HeapReleased, 10))
	report.AddUnConverted(metrics.TypeGauge, "HeapSys", strconv.FormatUint(memStats.HeapSys, 10))
	report.AddUnConverted(metrics.TypeGauge, "LastGC", strconv.FormatUint(memStats.LastGC, 10))
	report.AddUnConverted(metrics.TypeGauge, "Lookups", strconv.FormatUint(memStats.Lookups, 10))
	report.AddUnConverted(metrics.TypeGauge, "MCacheInuse", strconv.FormatUint(memStats.MCacheInuse, 10))
	report.AddUnConverted(metrics.TypeGauge, "MCacheSys", strconv.FormatUint(memStats.MCacheSys, 10))
	report.AddUnConverted(metrics.TypeGauge, "MSpanInuse", strconv.FormatUint(memStats.MSpanInuse, 10))
	report.AddUnConverted(metrics.TypeGauge, "MSpanSys", strconv.FormatUint(memStats.MSpanSys, 10))
	report.AddUnConverted(metrics.TypeGauge, "Mallocs", strconv.FormatUint(memStats.Mallocs, 10))
	report.AddUnConverted(metrics.TypeGauge, "NextGC", strconv.FormatUint(memStats.NextGC, 10))
	report.AddUnConverted(metrics.TypeGauge, "NumForcedGC", strconv.FormatUint(uint64(memStats.NumForcedGC), 10))
	report.AddUnConverted(metrics.TypeGauge, "NumGC", strconv.FormatUint(uint64(memStats.NumGC), 10))
	report.AddUnConverted(metrics.TypeGauge, "OtherSys", strconv.FormatUint(memStats.OtherSys, 10))
	report.AddUnConverted(metrics.TypeGauge, "PauseTotalNs", strconv.FormatUint(memStats.PauseTotalNs, 10))
	report.AddUnConverted(metrics.TypeGauge, "StackInuse", strconv.FormatUint(memStats.StackInuse, 10))
	report.AddUnConverted(metrics.TypeGauge, "StackSys", strconv.FormatUint(memStats.StackSys, 10))
	report.AddUnConverted(metrics.TypeGauge, "Sys", strconv.FormatUint(memStats.Sys, 10))
	report.AddUnConverted(metrics.TypeGauge, "TotalAlloc", strconv.FormatUint(memStats.TotalAlloc, 10))
	report.AddUnConverted(metrics.TypeGauge, "RandomValue", strconv.FormatUint(rand.Uint64(), 10))

	return report
}
