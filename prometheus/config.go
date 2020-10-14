package prometheus

import "github.com/spf13/viper"

type Config struct {
	Enable bool

	// The port to bind HTTP metrics API
	MetricsPort int
}

func InitConfig() (*Config, error) {
	config := &Config{
		Enable:      viper.GetBool("Prometheus.Enable"),
		MetricsPort: viper.GetInt("Prometheus.MetricsHTTPPort"),
	}

	if config.MetricsPort == 0 {
		config.MetricsPort = 24000
	}

	return config, nil
}
