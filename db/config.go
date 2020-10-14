package db

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUsername   string
	DBPassword   string
	DBName       string
	MaxOpenConns int32
}

func InitConfig() (*Config, error) {
	config := &Config{
		DBHost:       viper.GetString("PostgreSQL.DBHost"),
		DBPort:       viper.GetString("PostgreSQL.DBPort"),
		DBUsername:   viper.GetString("PostgreSQL.DBUsername"),
		DBPassword:   viper.GetString("PostgreSQL.DBPassword"),
		DBName:       viper.GetString("PostgreSQL.DBName"),
		MaxOpenConns: viper.GetInt32("PostgreSQL.MaxOpenConns"),
	}
	if config.DBHost == "" {
		config.DBHost = "localhost"
	}
	if config.DBPort == "" {
		config.DBPort = "5432"
	}
	if config.DBUsername == "" {
		config.DBUsername = "postgres"
	}
	if config.DBPassword == "" {
		config.DBPassword = "postgres"
	}
	if config.DBName == "" {
		config.DBName = "ticket_reservation"
	}
	return config, nil
}
