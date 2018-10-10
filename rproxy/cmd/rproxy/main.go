package main

import (
	"github.com/ggriffiths/gofun/rproxy/pkg/proxy"
	log "github.com/inconshreveable/log15"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("CACHE_EXPIRY", "60m")
	viper.SetDefault("CACHE_CAPACITY", 500)
	viper.SetDefault("STORAGE_URI", "localhost:6379")
	viper.SetDefault("WORKER_COUNT", 30)
}

func main() {
	log.Info("Starting rproxy!")

	viper.AutomaticEnv()
	if err := proxy.Start(&proxy.Config{
		Port:          viper.GetString("PORT"),
		CacheExpiry:   viper.GetDuration("CACHE_EXPIRY"),
		CacheCapacity: viper.GetInt64("CACHE_CAPACITY"),
		StorageURI:    viper.GetString("REDIS_URI"),
		WorkerCount:   viper.GetInt("WORKER_COUNT"),
	}); err != nil {
		log.Error("Proxy server failed", "err", err)
	}
}
