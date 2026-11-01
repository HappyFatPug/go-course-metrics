package main

import (
	"flag"

	"github.com/HappyFatPug/go-course-metrics/pkg/flags"
)

var config struct {
	netAddr        flags.NetAddress
	reportInterval int64
	pollInterval   int64
	debug          bool
}

func init() {
	config.netAddr.Host = "localhost"
	config.netAddr.Port = 8080
}

func parseFlags() {
	flag.BoolVar(&config.debug, "debug", false, "Debug mode")
	flag.Var(&config.netAddr, "a", "Net address host:port")
	flag.Int64Var(&config.pollInterval, "p", 2, "Poll interval in sec")
	flag.Int64Var(&config.reportInterval, "r", 10, "Report interval in sec")

	flag.Parse()
}
