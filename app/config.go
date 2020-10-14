package app

import "github.com/spf13/viper"

type Config struct {
	TokenSignerPrivateKeyPath string
	TokenSignerPublicKeyPath  string
}

func InitConfig() (*Config, error) {
	config := &Config{
		TokenSignerPrivateKeyPath: viper.GetString("TokenSignerPrivateKeyPath"),
		TokenSignerPublicKeyPath:  viper.GetString("TokenSignerPublicKeyPath"),
	}

	return config, nil
}
