package main

import (
	"fmt"
	"github.com/mixailo/go-training-metrics/internal/service/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"

	"github.com/mixailo/go-training-metrics/internal/service/metrics"
	"github.com/mixailo/go-training-metrics/internal/service/poller"
	"github.com/mixailo/go-training-metrics/internal/service/sender"
)

var totalPolls int64
var report metrics.Report

func gracefulShutdown(done chan bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		done <- true
	}()
}

func init() {
	logger.Initialize("info")
	logger.Log.Info("Starting agent")
}

func main() {
	agentConf := initConfig()
	reportEndpoint := sender.NewServerEndpoint("http", agentConf.endpoint.Host, agentConf.endpoint.Port)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	done := make(chan bool, 1)

	go gracefulShutdown(done)

	currentTime := time.Now()
	lastPoll := currentTime
	lastReport := currentTime
	for {
		select {
		case <-done:
			fmt.Println("Done, shutting down gracefully!")
			return
		case _ = <-ticker.C:
			if time.Since(lastPoll) >= time.Duration(agentConf.pollInterval)*time.Second {
				report = poller.PollMetrics()
				totalPolls += 1
				lastPoll = time.Now()
			}
			if time.Since(lastReport) >= time.Duration(agentConf.reportInterval)*time.Second {
				report.Add(metrics.Metrics{ID: "PollCount", MType: metrics.TypeCounter.String(), Delta: &totalPolls})
				err := sender.SendReport(report, reportEndpoint)
				if err != nil {
					logger.Log.Info("error", zap.Error(err))
				}
				lastReport = time.Now()
				totalPolls = 0
			}
		}
	}
}
