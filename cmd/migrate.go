package cmd

import (
	"fmt"
	"sort"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"ticket-reservation/db/migrations"
	"ticket-reservation/log"
)

var migrateDbCmd = &cobra.Command{
	Use: "migrate-db",
	RunE: func(cmd *cobra.Command, args []string) error {
		number, _ := cmd.Flags().GetInt("number")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		dbHost := viper.GetString("PostgreSQL.DBHost")
		dbPort := viper.GetString("PostgreSQL.DBPort")
		dbUser := viper.GetString("PostgreSQL.DBUser")
		if dbUser == "" {
			dbUser = "postgres"
		}
		dbPassword := viper.GetString("PostgreSQL.DBPassword")
		if dbPassword == "" {
			dbPassword = "postgres"
		}
		dbName := viper.GetString("PostgreSQL.DBName")

		logColor := viper.GetBool("Log.Color")
		logJSON := viper.GetBool("Log.JSON")
		logger, err := log.NewLogger(&log.Configuration{
			EnableConsole:     true,
			ConsoleLevel:      log.Debug,
			ConsoleJSONFormat: logJSON,
			Color:             logColor,
		}, log.InstanceZapLogger)
		if err != nil {
			return err
		}

		if dryRun {
			logger.Infof("=== DRY RUN ===")
		}

		sort.Slice(migrations.Migrations, func(i, j int) bool {
			return migrations.Migrations[i].Number < migrations.Migrations[j].Number
		})

		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost,
			dbPort,
			dbUser,
			dbPassword,
			dbName,
		)
		db, err := gorm.Open("postgres", connStr)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		defer db.Close()

		// Make sure Migration table is there
		logger.Debugf("ensuring migrations table is present")
		if err := db.AutoMigrate(&migrations.Migration{}).Error; err != nil {
			return errors.Wrap(err, "unable to automatically migrate migrations table")
		}

		var latest migrations.Migration
		if err := db.Order("number desc").First(&latest).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
			return errors.Wrap(err, "unable to find latest migration")
		}

		noMigrationsApplied := latest.Number == 0

		if noMigrationsApplied && len(migrations.Migrations) == 0 {
			logger.Infof("no migrations to apply")
			return nil
		}

		if latest.Number >= migrations.Migrations[len(migrations.Migrations)-1].Number {
			logger.Infof("no migrations to apply")
			return nil
		}

		if number == -1 {
			number = int(migrations.Migrations[len(migrations.Migrations)-1].Number)
		}

		if uint(number) <= latest.Number && latest.Number > 0 {
			logger.Infof("no migrations to apply, specified number is less than or equal to latest migration; backwards migrations are not supported")
			return nil
		}

		for _, migration := range migrations.Migrations {
			if migration.Number > uint(number) {
				break
			}

			if migration.Number <= latest.Number {
				continue
			}

			if latest.Number > 0 {
				logger.Infof("continuing migration starting from %d", migration.Number)
			}

			logger := logger.WithFields(log.Fields{
				"migration_number": migration.Number,
			})
			logger.Infof("applying migration %q", migration.Name)

			if dryRun {
				continue
			}

			tx := db.Begin()

			if err := migration.Forwards(tx); err != nil {
				logger.Errorf("unable to apply migration, rolling back. err: %+v", err)
				if err := tx.Rollback().Error; err != nil {
					logger.Errorf("unable to rollback... err: %+v", err)
				}
				break
			}

			if err := tx.Commit().Error; err != nil {
				logger.Errorf("unable to commit transaction... err: %+v", err)
				break
			}

			// Create migration record
			if err := db.Create(migration).Error; err != nil {
				logger.Errorf("unable to create migration record. err: %+v", err)
				break
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateDbCmd)

	migrateDbCmd.Flags().Int("number", -1, "the migration to run forwards until; if not set, will run all migrations")
	migrateDbCmd.Flags().Bool("dry-run", false, "print out migrations to be applied without running them")
}
