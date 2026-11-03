package config

import (
	"flag"

	"github.com/HappyFatPug/go-course-metrics/pkg/flags"
)

type ServerConfig struct {
	NetAddr flags.NetAddress
	Debug   bool
}

func ParseServerFlags() *ServerConfig {
	config := ServerConfig{
		NetAddr: flags.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
	}
	flag.BoolVar(&config.Debug, "debug", false, "Debug mode")
	flag.Var(&config.NetAddr, "a", "Net address host:port")

	flag.Parse()

	return &config
}

type AgentClientConfig struct {
	Scheme string
	NetAddr flags.NetAddress
}

type AgentConfig struct {
	ClientConfig   *AgentClientConfig
	ReportInterval int64
	PollInterval   int64
	Debug          bool
}

func ParseAgentFlags() *AgentConfig {
	clientConfig := AgentClientConfig{
		Scheme: "http",
		NetAddr: flags.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
	}
	config := AgentConfig{
		ClientConfig: &clientConfig,
	}

	flag.BoolVar(&config.Debug, "debug", false, "Debug mode")
	flag.Var(&clientConfig.NetAddr, "a", "Net address host:port")
	flag.Int64Var(&config.PollInterval, "p", 2, "Poll interval in sec")
	flag.Int64Var(&config.ReportInterval, "r", 10, "Report interval in sec")

	flag.Parse()

	return &config
}
