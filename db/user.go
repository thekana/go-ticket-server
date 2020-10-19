package db

import (
	"context"
	model "ticket-reservation/db/model"

	"github.com/pkg/errors"
)

type DBUserInterface interface {
	CreateUser(username string) (int64, error)
	GetUserById(id int64) (*model.UserWithRoleList, error)
	AssignRoleToUser(id int64) (int64, error) // FIXME: Not sure what to return here

	// TODO: [Phase 1] decide what to store in memory(pseudo-db) and to store in real db

	// admin/cust view all event
	// org create event [once created store quota in memory]
	// org view own event
	// org edit own event [update the final quota in memory]
	// org delete own event -> should also go and delete related reservations
	// get Event detail by ID for booking purpose
	// cust delete their reservations from table [in memory for now]
	// org fetch total ticket reserved / remaining (assumed for each event) return a list of such thing
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

func (pgdb *PostgresqlDB) GetUserById(id int64) (*model.UserWithRoleList, error) {
	// TODO: Just return user with role list
	userWithRole := &model.UserWithRoleList{}
	rows, err := pgdb.DB.Query(context.Background(), `
		select u.username, r.role from users u
		inner join user_roles ur on ur.user_id = u.id
		inner join roles r on r.id = ur.role_id
		where u.id = $1
		`, id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	// TODO: Loop rows and parse data properly check morchana/db/back_office
	return userWithRole, nil
}

// Requirement doesnt say anything about assigning roles to new user created via API
// Assumption: new users to always have customer role

func (pgdb *PostgresqlDB) AssignRoleToUser(id int64) (int64, error) {
	// TODO: FIXME: Not sure what to return here
	return 0, nil
}
