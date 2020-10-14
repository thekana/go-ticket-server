package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"ticket-reservation/db"
)

var createAdminUserCmd = &cobra.Command{
	Use:   "create-admin-user",
	Short: "Create admin user",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")

		if username == "" {
			return errors.New("username cannot be empty")
		}

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

		// TODO

		logger.Infof("Admin user created")

		return nil
	},
}

func init() {
	createAdminUserCmd.Flags().StringP("username", "u", "", "Username")
	createAdminUserCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(createAdminUserCmd)
}
