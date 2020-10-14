package db

import (
	"context"

	"github.com/pkg/errors"
)

type DBUserInterface interface {
	CreateUser(username string) (int64, error)
}

func (pgdb *PostgresqlDB) CreateUser(username string) (int64, error) {
	var userID int64

	err := pgdb.DB.QueryRow(context.Background(), `
		INSERT INTO users (
			"username"
		)
		VALUES ($1)
		RETURNING id
	`,
		username,
	).Scan(&userID)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to create user")
	}

	return userID, nil
}
