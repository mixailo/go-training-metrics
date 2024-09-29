package main

import (
	"github.com/mixailo/go-training-metrics/internal/metrics"
	"github.com/mixailo/go-training-metrics/internal/poller"
	"github.com/mixailo/go-training-metrics/internal/reporter"
	"log"
	"strconv"
	"time"
)

const pollInterval = 2

const reportInterval = 10

func main() {
	lastPoll := time.Now()
	lastReport := time.Now()
	var totalPolls int64
	var report metrics.Report
	reportEndpoint := reporter.NewServerEndpoint("http", "127.0.0.1", 8080)

	for {
		currentTime := time.Now()
		if (currentTime.Sub(lastPoll).Seconds()) >= pollInterval {
			report = poller.PollMetrics()
			totalPolls += 1

			lastPoll = currentTime
		}
		if currentTime.Sub(lastReport).Seconds() >= reportInterval {
			report.Add(metrics.TypeCounter, "PollCount", strconv.FormatInt(totalPolls, 10))

			err := reporter.SendReport(report, reportEndpoint)
			if err != nil {
				log.Print(err.Error())
			}
			lastReport = currentTime
			totalPolls = 0
		}
		time.Sleep(500 * time.Millisecond)
	}
}
