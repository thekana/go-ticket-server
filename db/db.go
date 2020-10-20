package db

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"ticket-reservation/log"
)

type DB interface {
	DBUserInterface
	DBEventInterface
	//DBReservationInterface
	Close() error
	PrintSystem()
}

type PostgresqlDB struct {
	logger   log.Logger
	DB       *pgxpool.Pool
	MemoryDB *System
}

func New(config *Config, logger log.Logger) (pgdb *PostgresqlDB, err error) {
	pgdb = &PostgresqlDB{
		logger: logger.WithFields(log.Fields{
			"module": "db",
		}),
		MemoryDB: NewSystem(), // Init memoryDB
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUsername,
		config.DBPassword,
		config.DBName,
	)

	//db, err = pgx.Connect(context.Background(), connStr)
	var connectConf, _ = pgxpool.ParseConfig(connStr)
	connectConf.MaxConns = config.MaxOpenConns
	//connectConf.MaxConnLifetime = 300 * time.Second // use defaults until we have benchmarked this further
	//connectConf.HealthCheckPeriod = 300 * time.Second
	//connectConf.ConnConfig.PreferSimpleProtocol = true // don't wrap queries into transactions
	connectConf.ConnConfig.Logger = NewDatabaseLogger(&pgdb.logger)
	connectConf.ConnConfig.LogLevel = pgx.LogLevelWarn
	pgdb.DB, err = pgxpool.ConnectConfig(context.Background(), connectConf)
	if err != nil {
		pgdb.logger.Errorf("Error connecting to postgres: %+v")
		return nil, err
	}

	return pgdb, nil
}

func (pgdb *PostgresqlDB) Close() error {
	pgdb.DB.Close()
	return nil
}

func (pgdb *PostgresqlDB) PrintSystem() {
	spew.Dump(pgdb.MemoryDB)
}
