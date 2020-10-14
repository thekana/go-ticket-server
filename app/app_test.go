package app

import (
	"github.com/spf13/viper"

	log "ticket-reservation/log"
)

func getLoggerForTesting() (log.Logger, error) {
	logLevel := viper.GetString("Log.Level")
	logLevel = log.NormalizeLogLevel(logLevel)

	logColor := viper.GetBool("Log.Color")
	logJSON := viper.GetBool("Log.JSON")

	logger, err := log.NewLogger(&log.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logLevel,
		ConsoleJSONFormat: logJSON,
		Color:             logColor,
	}, log.InstanceZapLogger)
	if err != nil {
		return nil, err
	}
	return logger, nil
}
