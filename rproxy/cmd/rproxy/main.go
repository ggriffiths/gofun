package main

import (
	"os"

	"github.com/ggriffiths/gofun/rproxy/pkg/proxy"
	log "github.com/inconshreveable/log15"
	"github.com/kelseyhightower/envconfig"
)

const (
	msgFailedEnv = "Failed to get and parse environment variable"
)

func main() {
	log.Info("Starting rproxy!")

	// Get env variables
	var cfg proxy.Config
	err := envconfig.Process("RProxyServer", &cfg)
	if err != nil {
		log.Crit("Failed to parse enviroment variables", "err", err)
		os.Exit(1)
	}

	// Create and start proxy server
	s := proxy.NewServer(&cfg)
	if err := s.Start(); err != nil {
		log.Error("ProxyServer server failed", "err", err)
		s.Close()
		os.Exit(1)
	}
	defer s.Close()
}
