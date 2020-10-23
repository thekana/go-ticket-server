package db

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type DBUserInterface interface {
	CreateUser(username string) (int64, error)
	GetUserById(id int64) (*model.UserWithRoleList, error)
	AssignRoleToUser(id int64, role model.Role) (int64, error)
	GetUserByName(name string) (*model.UserWithRoleList, error)
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
	pgdb.MemoryDB.AddUserToSystem(NewUserData(username, int(userID)))
	//fmt.Printf("-----------------\n")
	//spew.Dump(pgdb.MemoryDB)
	//fmt.Printf("-----------------\n")
	return userID, nil
}

func (pgdb *PostgresqlDB) GetUserById(id int64) (*model.UserWithRoleList, error) {
	userWithRole := &model.UserWithRoleList{}
	rows, err := pgdb.DB.Query(context.Background(), `
		select u.id as uid ,u.username, r.role from users u
		inner join user_roles ur on ur.user_id = u.id
		inner join roles r on r.id = ur.role_id
		where u.id = $1
		`, id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var role string
		var id int64
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

func (pgdb *PostgresqlDB) AssignRoleToUser(id int64, role model.Role) (int64, error) {
	var rowId int64
	err := pgdb.DB.QueryRow(context.Background(), `
		INSERT INTO user_roles(user_id,role_id) values ($1,$2)
		RETURNING id
		`, id, role).Scan(&rowId)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to create user")
	}
	return rowId, nil
}
