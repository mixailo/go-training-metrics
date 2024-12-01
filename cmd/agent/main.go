package main

import (
	"github.com/mixailo/go-training-metrics/internal/service/logger"
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

var agentConf config

func main() {
	agentConf = initConfig()
	if err := logger.Initialize(agentConf.logLevel); err != nil {
		panic(err)
	}
	logger.Log.Info("agent start")
	gracefulShutdown()

	lastPoll := time.Now()
	lastReport := lastPoll

	reportEndpoint := sender.NewServerEndpoint("http", agentConf.endpoint.Host, agentConf.endpoint.Port)
	for {
		time.Sleep(100 * time.Millisecond)
		currentTime := time.Now()
		if (currentTime.Sub(lastPoll).Milliseconds()) >= agentConf.pollInterval*1000 {
			report = poller.PollMetrics()
			totalPolls += 1
			lastPoll = currentTime
		}
		if currentTime.Sub(lastReport).Milliseconds() >= agentConf.reportInterval*1000 {
			report.Add(metrics.Metrics{ID: "PollCount", MType: metrics.TypeCounter.String(), Delta: &totalPolls})
			err := sender.SendReport(report, reportEndpoint)
			if err != nil {
				log.Print(err.Error())
			}
			lastReport = currentTime
			totalPolls = 0
		}

	}
}
