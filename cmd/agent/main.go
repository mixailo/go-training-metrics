package main

import (
	"fmt"
	"github.com/mixailo/go-training-metrics/internal/service/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mixailo/go-training-metrics/internal/service/metrics"
	"github.com/mixailo/go-training-metrics/internal/service/poller"
	"github.com/mixailo/go-training-metrics/internal/service/sender"
)

var totalPolls int64
var report metrics.Report

// yet ungraceful
func gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Shutting down gracefully")
		os.Exit(0)
	}()
}

func main() {
	agentConf := initConfig()
	gracefulShutdown()

	logger.Initialize("info")

	logger.Log.Info("Starting agent")
	logger.Log.Info(fmt.Sprint(agentConf))
	lastPoll := time.Now()
	lastReport := time.Now()

	reportEndpoint := sender.NewServerEndpoint("http", agentConf.endpoint.Host, agentConf.endpoint.Port)
	time.Sleep(2 * time.Second)
	for {
		currentTime := time.Now()
		if time.Since(lastPoll) >= time.Duration(agentConf.pollInterval)*time.Second {
			report = poller.PollMetrics()
			totalPolls += 1
			lastPoll = currentTime
		}
		if time.Since(lastReport) >= time.Duration(agentConf.reportInterval)*time.Second {
			report.Add(metrics.Metrics{ID: "PollCount", MType: metrics.TypeCounter.String(), Delta: &totalPolls})
			err := sender.SendReport(report, reportEndpoint)
			if err != nil {
				logger.Log.Info("error", zap.Error(err))
			}
			lastReport = currentTime
			totalPolls = 0
		}
		time.Sleep(100 * time.Millisecond)
	}
}
