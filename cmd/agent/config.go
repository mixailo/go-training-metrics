package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type endpoint struct {
	Host string
	Port int
}

func (e *endpoint) parse(value string) error {
	vals := strings.Split(value, ":")

	if len(vals) != 2 {
		return errors.New("endpoint string must have exactly one divider")
	}

	e.Host = vals[0]
	port, err := strconv.Atoi(vals[1])
	if err != nil {

		return err
	}
	e.Port = port

	return nil
}

type config struct {
	endpoint       endpoint
	pollInterval   int64
	reportInterval int64
	logLevel       string
}

func (e *endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func (e *endpoint) Set(value string) error {
	items := strings.Split(value, ":")
	if len(items) != 2 {
		return errors.New("invalid format")
	}

	e.Host = items[0]
	e.Port, _ = strconv.Atoi(items[1])

	return nil
}

func defaultConfig() (cfg config) {
	cfg = config{
		endpoint: endpoint{
			Host: "localhost",
			Port: 8080,
		},
		pollInterval:   2,
		reportInterval: 10,
		logLevel:       "info",
	}
	return
}

func envConfig(defCfg config) (cfg config) {
	cfg = defCfg

	v, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		interval, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		cfg.reportInterval = interval
	}

	v, ok = os.LookupEnv("POLL_INTERVAL")
	if ok {
		interval, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		cfg.pollInterval = interval
	}

	v, ok = os.LookupEnv("ADDRESS")
	if ok {
		// parse env vars
		err := cfg.endpoint.parse(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	v, ok = os.LookupEnv("LOG_LEVEL")
	if ok {
		cfg.logLevel = v
	}

	return cfg
}

func initConfig() config {
	return argsConfig(envConfig(defaultConfig()))
}

func argsConfig(cfg config) config {
	flag.Var(&cfg.endpoint, "a", "server endpoint")
	flag.Int64Var(&cfg.pollInterval, "p", cfg.pollInterval, "poll interval")
	flag.Int64Var(&cfg.reportInterval, "r", cfg.reportInterval, "report interval")
	flag.StringVar(&cfg.logLevel, "l", cfg.logLevel, "log level [info]")

	flag.Parse()
	return cfg
}
