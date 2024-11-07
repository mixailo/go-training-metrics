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
	endpoint        endpoint
	logLevel        string
	storeInterval   int64
	fileStoragePath string
	doRestoreValues bool
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

func (c *config) validate() error {
	if c.storeInterval < 0 {
		return errors.New("store interval must be a positive number or zero")
	}

	return nil
}

func initConfig() (config, error) {
	c := envConfig(argsConfig(defaultConfig()))
	err := c.validate()

	return c, err
}

func envConfig(defCfg config) config {
	var (
		v  string
		ok bool
	)

	cfg := defCfg
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

	v, ok = os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		cfg.fileStoragePath = v
	}

	v, ok = os.LookupEnv("STORE_INTERVAL")
	if ok {
		vv, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			cfg.storeInterval = vv
		}
	}

	v, ok = os.LookupEnv("RESTORE")
	if ok {
		cfg.doRestoreValues = v == "true"
	}

	return cfg
}

func defaultConfig() (cfg config) {
	return config{
		endpoint: endpoint{
			host: "localhost",
			port: 8080,
		},
		logLevel:        "info",
		doRestoreValues: true,
		storeInterval:   300,
		fileStoragePath: "values.json",
	}
}

func argsConfig(cfg config) config {
	flag.Var(&cfg.endpoint, "a", "server endpoint [host:port]")
	flag.StringVar(&cfg.logLevel, "l", "info", "log level [info]")
	flag.BoolVar(&cfg.doRestoreValues, "r", false, "do restore saved values")
	flag.StringVar(&cfg.fileStoragePath, "f", cfg.fileStoragePath, "path to storage file")
	flag.Int64Var(&cfg.storeInterval, "i", cfg.storeInterval, "storage save interval in seconds")
	flag.Parse()

	return cfg
}
