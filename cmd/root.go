package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"ticket-reservation/log"
)

var rootCmd = &cobra.Command{
	Use:   "ticket-reservation",
	Short: "Ticker Reservation",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var configFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("Log.Level", "debug")
	viper.SetDefault("Log.Color", true)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("unable to read config: %v\n", err)
		os.Exit(1)
	}
}

func getLogger() (log.Logger, error) {
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
