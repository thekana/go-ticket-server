package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"ticket-reservation/db/model"
	"ticket-reservation/log"
)

type DB interface {
	DBUserInterface
	DBEventInterface
	DBReservationInterface
	Close() error
	PopulateSystem()
	SeedDBForTest()
}

type PostgresqlDB struct {
	logger log.Logger
	DB     *pgxpool.Pool
}

func New(config *Config, logger log.Logger) (pgdb *PostgresqlDB, err error) {
	pgdb = &PostgresqlDB{
		logger: logger.WithFields(log.Fields{
			"module": "db",
		}),
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

func (pgdb *PostgresqlDB) PopulateSystem() {
	// Create 1 admin
	_, _ = pgdb.CreateUser("admin", model.Admin)
	// Create 2 orgs
	org1ID, _ := pgdb.CreateUser("org1", model.Organizer)
	org2ID, _ := pgdb.CreateUser("org2", model.Organizer)
	// Create 1 cust
	_, _ = pgdb.CreateUser("cust1", model.Customer)
	_, _ = pgdb.CreateUser("cust2", model.Customer)
	_, _ = pgdb.CreateUser("cust3", model.Customer)
	// Each org create two events
	_, _ = pgdb.CreateEvent(org1ID, "org1 event1", 10000)
	_, _ = pgdb.CreateEvent(org1ID, "org1 event2", 10000)
	// Each org create two events
	_, _ = pgdb.CreateEvent(org2ID, "org2 event1", 10000)
	_, _ = pgdb.CreateEvent(org2ID, "org2 event2", 10000)
}

func (pgdb *PostgresqlDB) SeedDBForTest() {
	// Create 1 admin
	_, _ = pgdb.CreateUser("admin", model.Admin)
	// Create 2 orgs
	org1ID, _ := pgdb.CreateUser("org1", model.Organizer)
	org2ID, _ := pgdb.CreateUser("org2", model.Organizer)
	// Create 1 cust
	_, _ = pgdb.CreateUser("cust1", model.Customer)
	_, _ = pgdb.CreateUser("cust2", model.Customer)
	_, _ = pgdb.CreateUser("cust3", model.Customer)
	// Each org create two events
	_, _ = pgdb.CreateEvent(org1ID, "org1 event1", 10000)
	_, _ = pgdb.CreateEvent(org1ID, "org1 event2", 10000)
	// Each org create two eents
	_, _ = pgdb.CreateEvent(org2ID, "org2 event1", 10000)
	_, _ = pgdb.CreateEvent(org2ID, "org2 event2", 10000)
}
