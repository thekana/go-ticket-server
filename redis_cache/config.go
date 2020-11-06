package redis_cache

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	RedisMode              string
	RedisHost              string
	RedisPort              string
	SentinelAddrs          []string
	SentinelMasterName     string
	RedisPassword          string
	RedisDB                int
	MaxRetries             int
	MinRetryBackoffSeconds int
	MaxRetryBackoffSeconds int
	DialTimeoutSeconds     int
	WriteTimeoutSeconds    int
	PoolTimeoutSeconds     int
	RedisConnectionTimeout time.Duration
}

func InitConfig() (*Config, error) {
	config := &Config{
		RedisMode:              viper.GetString("RedisCache.RedisMode"),
		RedisHost:              viper.GetString("RedisCache.RedisHost"),
		RedisPort:              viper.GetString("RedisCache.RedisPort"),
		SentinelMasterName:     viper.GetString("RedisCache.Sentinel.MasterName"),
		SentinelAddrs:          viper.GetStringSlice("RedisCache.Sentinel.Addrs"),
		RedisPassword:          viper.GetString("RedisCache.RedisPassword"),
		RedisDB:                viper.GetInt("RedisCache.RedisDB"),
		MaxRetries:             viper.GetInt("RedisCache.MaxRetries"),
		MinRetryBackoffSeconds: viper.GetInt("RedisCache.MinRetryBackoffSeconds"),
		MaxRetryBackoffSeconds: viper.GetInt("RedisCache.MaxRetryBackoffSeconds"),
		DialTimeoutSeconds:     viper.GetInt("RedisCache.DialTimeoutSeconds"),
		WriteTimeoutSeconds:    viper.GetInt("RedisCache.WriteTimeoutSeconds"),
		PoolTimeoutSeconds:     viper.GetInt("RedisCache.PoolTimeoutSeconds"),
		RedisConnectionTimeout: viper.GetDuration("RedisCache.RedisConnectionTimeout"),
	}
	if config.RedisHost == "" {
		config.RedisHost = "localhost"
	}
	if config.RedisPort == "" {
		config.RedisPort = "6379"
	}
	return config, nil
}
