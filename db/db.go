package db

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"ticket-reservation/db/model"

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
	PopulateSystem()
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

func (pgdb *PostgresqlDB) PopulateSystem() {
	// Create 1 admin
	adminID, _ := pgdb.CreateUser("admin")
	_, _ = pgdb.AssignRoleToUser(adminID, model.Admin)
	// Create 2 orgs
	org1ID, _ := pgdb.CreateUser("org1")
	org2ID, _ := pgdb.CreateUser("org2")
	_, _ = pgdb.AssignRoleToUser(org1ID, model.Organizer)
	_, _ = pgdb.AssignRoleToUser(org2ID, model.Organizer)
	// Create 1 cust
	cust1ID, _ := pgdb.CreateUser("cust1")
	_, _ = pgdb.AssignRoleToUser(cust1ID, model.Organizer)
	// Each org create two events
	_, _ = pgdb.CreateEvent(int(org1ID), "org1 event1", 1000)
	_, _ = pgdb.CreateEvent(int(org1ID), "org1 event2", 1000)
	// Each org create two events
	_, _ = pgdb.CreateEvent(int(org2ID), "org2 event1", 1000)
	_, _ = pgdb.CreateEvent(int(org2ID), "org2 event2", 1000)
}
