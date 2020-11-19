package cmd

import (
	"github.com/spf13/cobra"
	"ticket-reservation/db"
)

// Seed DB for testing purposes only

var createSeedCmd = &cobra.Command{
	Use:   "seed-db",
	Short: "Seed DB",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := getLogger()
		if err != nil {
			return err
		}

		dbConfig, err := db.InitConfig()
		if err != nil {
			return err
		}

		db, err := db.New(dbConfig, logger)
		if err != nil {
			return err
		}
		defer db.Close()

		// Create the user
		db.SeedDBForTest()

		logger.Infof("DB seeded for testing")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createSeedCmd)
}
