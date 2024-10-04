package main

import (
	"github.com/mixailo/go-training-metrics/internal/metrics"
	"github.com/mixailo/go-training-metrics/internal/poller"
	"github.com/mixailo/go-training-metrics/internal/reporter"
	"log"
	"strconv"
	"time"
)

var totalPolls int64
var report metrics.Report

func main() {
	ParseFlags()
	lastPoll := time.Now()
	lastReport := time.Now()

	reportEndpoint := reporter.NewServerEndpoint("http", serverEndpoint.host, serverEndpoint.port)

	for {
		currentTime := time.Now()
		if (currentTime.Sub(lastPoll).Milliseconds()) >= pollInterval*1000 {
			report = poller.PollMetrics()
			totalPolls += 1

			lastPoll = currentTime
		}
		if currentTime.Sub(lastReport).Milliseconds() >= reportInterval*1000 {
			report.Add(metrics.TypeCounter, "PollCount", strconv.FormatInt(totalPolls, 10))

			err := reporter.SendReport(report, reportEndpoint)
			if err != nil {
				log.Print(err.Error())
			}
			lastReport = currentTime
			totalPolls = 0
		}
		time.Sleep(100 * time.Millisecond)
	}
}
