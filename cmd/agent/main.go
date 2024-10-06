package main

import (
	"fmt"
	"github.com/mixailo/go-training-metrics/internal/service/metrics"
	"github.com/mixailo/go-training-metrics/internal/service/poller"
	"github.com/mixailo/go-training-metrics/internal/service/reporter"
	"log"
	"strconv"
	"time"
)

var totalPolls int64
var report metrics.Report

func main() {
	agentConf := InitConfig()

	fmt.Println(agentConf)

	lastPoll := time.Now()
	lastReport := time.Now()

	reportEndpoint := reporter.NewServerEndpoint("http", agentConf.Endpoint.Host, agentConf.Endpoint.Port)

	for {
		currentTime := time.Now()
		if (currentTime.Sub(lastPoll).Milliseconds()) >= agentConf.PollInterval*1000 {
			report = poller.PollMetrics()
			totalPolls += 1

			lastPoll = currentTime
		}
		if currentTime.Sub(lastReport).Milliseconds() >= agentConf.ReportInterval*1000 {
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
