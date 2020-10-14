package http_api

import "github.com/spf13/viper"

type Config struct {
	// The port to bind HTTP application API server to
	Port int

	// The number of proxies positioned in front of the API. This is used to interpret
	// X-Forwarded-For headers.
	ProxyCount int

	LogLevel string
}

func InitConfig() (*Config, error) {
	config := &Config{
		Port:     viper.GetInt("API.HTTPServerPort"),
		LogLevel: viper.GetString("Log.Level"),
	}
	if config.Port == 0 {
		config.Port = 9092
	}
	return config, nil
}
