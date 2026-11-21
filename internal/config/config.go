package config

import (
	"flag"
	"log"

	"github.com/HappyFatPug/go-course-metrics/pkg/flags"
	"github.com/caarlos0/env/v6"
)

type StoreConfig struct {
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

type ServerConfig struct {
	NetAddr  flags.NetAddress `env:"ADDRESS"`
	StoreCfg *StoreConfig
	Debug    bool
}

func GetServerConfig() *ServerConfig {
	config := ServerConfig{
		NetAddr: flags.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
		StoreCfg: &StoreConfig{},
	}
	flag.BoolVar(&config.Debug, "debug", false, "Debug mode")
	flag.Var(&config.NetAddr, "a", "Net address host:port")
	flag.Int64Var(&config.StoreCfg.StoreInterval, "i", 300, "Store interval")
	flag.StringVar(&config.StoreCfg.FileStoragePath, "f", "./dump.json", "File storeage path")
	flag.BoolVar(&config.StoreCfg.Restore, "r", false, "Restore from store file")

	// parse cli arguments
	flag.Parse()

	// parse env variables
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

type AgentClientConfig struct {
	Scheme  string
	NetAddr flags.NetAddress `env:"ADDRESS"`
	UseGzip bool
}

type AgentConfig struct {
	ClientConfig   *AgentClientConfig
	ReportInterval int64 `env:"REPORT_INTERVAL"`
	PollInterval   int64 `env:"POLL_INTERVAL"`
	Debug          bool
}

func GetAgentConfig() *AgentConfig {
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
	flag.BoolVar(&clientConfig.UseGzip, "gzip", true, "Gzip mode")
	flag.Var(&clientConfig.NetAddr, "a", "Net address host:port")
	flag.Int64Var(&config.PollInterval, "p", 2, "Poll interval in sec")
	flag.Int64Var(&config.ReportInterval, "r", 10, "Report interval in sec")

	// parse cli arguments
	flag.Parse()

	// parse env variables
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
