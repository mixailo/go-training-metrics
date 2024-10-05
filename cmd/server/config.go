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
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

type Config struct {
	Endpoint Endpoint
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

func InitConfig() Config {
	return ArgsConfig(EnvConfig(DefaultConfig()))
}

func EnvConfig(defCfg Config) Config {
	cfg := defCfg
	v, ok := os.LookupEnv("ADDRESS")

	if ok {
		// parse env vars
		err := cfg.Endpoint.Parse(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cfg
}

func DefaultConfig() (cfg Config) {
	return Config{
		Endpoint: Endpoint{
			Host: "localhost",
			Port: 8080,
		},
	}
}

func ArgsConfig(cfg Config) Config {
	flag.Var(&cfg.Endpoint, "a", "server endpoint [host:port]")
	flag.Parse()
	return cfg
}
