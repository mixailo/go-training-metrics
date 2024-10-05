package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

func (c *Config) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) Set(value string) error {
	items := strings.Split(value, ":")
	if len(items) != 2 {
		return errors.New("invalid format")
	}

	c.Host = items[0]
	c.Port, _ = strconv.Atoi(items[1])

	return nil
}

func ParseEnvConfig(value string) (c Config, err error) {
	vals := strings.Split(value, ":")

	if len(vals) != 2 {
		return c, errors.New("endpoint string must have exactly one divider")
	}

	c.Host = vals[0]
	c.Port, err = strconv.Atoi(vals[1])
	if err != nil {
		return c, errors.New("invalid port number, cannot be parsed")
	}

	return c, nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func InitConfig() (cfg Config) {
	envConf, envConfSet := EnvConfig()
	argsConf, argsConfSet := ArgsConfig()

	if isFlagPassed("f") && argsConfSet {
		return argsConf
	} else if envConfSet {
		return envConf
	} else if argsConfSet {
		return argsConf
	} else {
		return DefaultConfig()
	}
}

func EnvConfig() (cfg Config, ok bool) {
	v, ok := os.LookupEnv("ADDRESS")

	if ok {
		// parse env vars
		cfg, err := ParseEnvConfig(v)
		if err != nil {
			return cfg, false
		} else {
			return cfg, true
		}
	}

	return
}

func DefaultConfig() (cfg Config) {
	return Config{
		Host: "localhost",
		Port: 8080,
	}
}

func ArgsConfig() (Config, bool) {
	cfg := DefaultConfig()
	flag.Var(&cfg, "a", "server endpoint [host:port]")
	_ = flag.Bool("f", false, "prefer command-line arguments for configuration")
	flag.Parse()

	return cfg, isFlagPassed("a")
}
