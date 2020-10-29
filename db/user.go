package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type DBUserInterface interface {
	CreateUser(username string, role model.Role) (int, error)
	GetUserById(id int) (*model.UserWithRoleList, error)
	GetUserByName(name string) (*model.UserWithRoleList, error)
}

func (pgdb *PostgresqlDB) CreateUser(username string, role model.Role) (int, error) {
	var userID int
	tx, err := pgdb.DB.Begin(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(context.Background())
		}
	}()
	err = tx.QueryRow(context.Background(), `
		INSERT INTO users (
			"username"
		)
		VALUES ($1)
		RETURNING id
	`,
		username,
	).Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return 0, errors.New("Duplicate username")
			}
		}
		return 0, errors.Wrap(err, "Unable to create user")
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO user_roles(user_id,role_id) values ($1,$2)`, userID, role)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to add user role on create")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "Unable to commit a transaction")
	}
	return userID, nil
}

func (pgdb *PostgresqlDB) GetUserById(id int) (*model.UserWithRoleList, error) {
	userWithRole := &model.UserWithRoleList{}
	rows, err := pgdb.DB.Query(context.Background(), `
		select u.id as uid ,u.username, r.role from users u
		inner join user_roles ur on ur.user_id = u.id
		inner join roles r on r.id = ur.role_id
		where u.id = $1
		`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role string
		var username string
		err = rows.Scan(nil, &username, &role)
		if err != nil {
			return nil, err
		}
		userWithRole.RoleList = append(userWithRole.RoleList, role)
		userWithRole.Username = username
	}
	userWithRole.ID = id
	return userWithRole, nil
}

func (pgdb *PostgresqlDB) GetUserByName(name string) (*model.UserWithRoleList, error) {
	userWithRole := &model.UserWithRoleList{}
	rows, err := pgdb.DB.Query(context.Background(), `
		select u.id as uid ,u.username, r.role from users u
		inner join user_roles ur on ur.user_id = u.id
		inner join roles r on r.id = ur.role_id
		where u.username = $1
		`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role string
		var id int
		err = rows.Scan(&id, nil, &role)
		if err != nil {
			return nil, err
		}
		userWithRole.RoleList = append(userWithRole.RoleList, role)
		userWithRole.ID = id
	}
	userWithRole.Username = name
	if len(userWithRole.RoleList) == 0 {
		return nil, &customError.UserError{
			Code:           customError.UserNotFound,
			Message:        "User not found",
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	return userWithRole, nil
}
