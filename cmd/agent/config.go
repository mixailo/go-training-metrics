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

type Endpoint struct {
	Host string
	Port int
}

func (e *Endpoint) Parse(value string) error {
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

type Config struct {
	Endpoint       Endpoint
	PollInterval   int64
	ReportInterval int64
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func (e *Endpoint) Set(value string) error {
	items := strings.Split(value, ":")
	if len(items) != 2 {
		return errors.New("invalid format")
	}

	e.Host = items[0]
	e.Port, _ = strconv.Atoi(items[1])

	return nil
}

func DefaultConfig() (cfg Config) {
	cfg = Config{
		Endpoint: Endpoint{
			Host: "localhost",
			Port: 8080,
		},
		PollInterval:   2,
		ReportInterval: 10,
	}
	return
}

func EnvConfig(defCfg Config) (cfg Config) {
	cfg = defCfg

	v, ok := os.LookupEnv("REPORT_INTERVAL")
	fmt.Println(v, ok)
	if ok {
		interval, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		cfg.ReportInterval = interval
	}

	v, ok = os.LookupEnv("POLL_INTERVAL")
	fmt.Println(v, ok)
	if ok {
		interval, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		cfg.PollInterval = interval
	}

	v, ok = os.LookupEnv("ADDRESS")
	fmt.Println(v, ok)
	if ok {
		// parse env vars
		err := cfg.Endpoint.Parse(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cfg
}

func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func InitConfig() Config {
	return ArgsConfig(EnvConfig(DefaultConfig()))
}

func ArgsConfig(cfg Config) Config {
	flag.Var(&cfg.Endpoint, "a", "server endpoint")
	flag.Int64Var(&cfg.PollInterval, "p", cfg.PollInterval, "poll interval")
	flag.Int64Var(&cfg.ReportInterval, "r", cfg.ReportInterval, "report interval")

	flag.Parse()
	return cfg
}
