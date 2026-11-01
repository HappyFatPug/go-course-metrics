package main

import (
	"flag"

	"github.com/HappyFatPug/go-course-metrics/pkg/flags"
)

var config struct {
	netAddr flags.NetAddress
	debug   bool
}

func init() {
	config.netAddr.Host = "localhost"
	config.netAddr.Port = 8080
}

func parseFlags() {
	flag.BoolVar(&config.debug, "debug", false, "Debug mode")
	flag.Var(&config.netAddr, "a", "Net address host:port")

	flag.Parse()
}
