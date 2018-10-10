package main

import (
	"github.com/ggriffiths/gofun/rproxy/pkg/proxy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("CACHE_EXPIRY", "60m")
	viper.SetDefault("CACHE_CAPACITY", 500)
	viper.SetDefault("STORAGE_URI", "localhost:6380")
}

func main() {
	viper.AutomaticEnv()

	if err := proxy.Start(&proxy.Config{
		Port:          viper.GetString("PORT"),
		CacheExpiry:   viper.GetDuration("CACHE_EXPIRY"),
		CacheCapacity: viper.GetInt64("CACHE_CAPACITY"),
		StorageURI:    viper.GetString("REDIS_URI"),
	}); err != nil {
		log.Error("Proxy server failed", "err", err)
	}
}
