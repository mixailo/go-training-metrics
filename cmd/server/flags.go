package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type endpoint struct {
	host string
	port int
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

var serverEndpoint = endpoint{host: "127.0.0.1", port: 8080}

func ParseFlags() {
	flag.Var(&serverEndpoint, "a", "server endpoint [host:port]")
	flag.Parse()
}
