package app

import (
	"errors"
)

// RegisterParams is
type RegisterParams struct {
	Username string `json:"username" validate:"required"`
}

// RegisterResult is
type RegisterResult struct {
	ID int64 `json:"id"`
}

// LoginParams is
type LoginParams struct {
	Username string `json:"username" validate:"required"`
}

// LoginResult is
type LoginResult struct {
	AuthToken string `json:"authToken"`
	UserID    int64  `json:"userId"`
}

// GetLoggedInInfoParams is just here to check if token is working
type GetLoggedInInfoParams struct {
	AuthToken string `json:"authToken" validate:"required"`
}

// GetLoggedInInfoResult is
type GetLoggedInInfoResult struct {
	UserID int64       `json:"userId"`
	Data   interface{} `json:"data"`
}

// GetLoggedInInfo checks validation for now
func (ctx *Context) GetLoggedInInfo(params GetLoggedInInfoParams) (*GetLoggedInInfoResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	_, claims, err := ctx.verifyToken(params.AuthToken)
	if err != nil {
		return nil, err
	}
	return &GetLoggedInInfoResult{
		UserID: 1,
		Data:   claims,
	}, nil
}

// FIXME: Temp solution
var (
	usernameMap map[string]int64 = make(map[string]int64)
	userCount   int64            = 0
)

// Login is a backend function
func (ctx *Context) Login(params LoginParams) (*LoginResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	// TODO: query database instead check db/user.go createUser
	id, ok := usernameMap[params.Username]
	if !ok {
		return nil, errors.New("User does not exist")
	}
	authToken, err := ctx.createToken(id)
	if err != nil {
		return nil, err
	}
	return &LoginResult{AuthToken: authToken, UserID: id}, nil
}

// Register is a backend function
func (ctx *Context) Register(params RegisterParams) (*RegisterResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	if _, ok := usernameMap[params.Username]; ok {
		return nil, errors.New("Username exists")
	}
	id := userCount
	usernameMap[params.Username] = id
	userCount++

	// TODO: Change to update database check db/user.go createUser

	return &RegisterResult{ID: id}, nil
}
