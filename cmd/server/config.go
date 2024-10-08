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
	host string `env:"HOST"`
	port int    `env:"PORT"`
}

type config struct {
	endpoint endpoint
}

func (e *endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.host, e.port)
}

func (e *endpoint) Set(value string) error {
	items := strings.Split(value, ":")
	if len(items) != 2 {
		return errors.New("invalid format")
	}

	e.host = items[0]
	e.port, _ = strconv.Atoi(items[1])

	return nil
}

func (e *endpoint) parse(value string) error {
	vals := strings.Split(value, ":")

	if len(vals) != 2 {
		return errors.New("endpoint string must have exactly one divider")
	}

	e.host = vals[0]
	port, err := strconv.Atoi(vals[1])
	if err != nil {

		return err
	}
	e.port = port

	return nil
}

func initConfig() config {
	return argsConfig(envConfig(defaultConfig()))
}

func envConfig(defCfg config) config {
	cfg := defCfg
	v, ok := os.LookupEnv("ADDRESS")

	if ok {
		// parse env vars
		err := cfg.endpoint.parse(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cfg
}

func defaultConfig() (cfg config) {
	return config{
		endpoint: endpoint{
			host: "localhost",
			port: 8080,
		},
	}
}

func argsConfig(cfg config) config {
	flag.Var(&cfg.endpoint, "a", "server endpoint [host:port]")
	flag.Parse()
	return cfg
}
