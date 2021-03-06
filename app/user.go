package app

import (
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

// RegisterParams is
type RegisterParams struct {
	Username string `json:"username" validate:"required"`
}

// RegisterResult is
type RegisterResult struct {
	ID int `json:"id"`
}

// LoginParams is
type LoginParams struct {
	Username string `json:"username" validate:"required"`
}

// LoginResult is
type LoginResult struct {
	AuthToken string `json:"authToken"`
	UserID    int    `json:"userId"`
	Username  string `json:"username"`
	Role      string `json:"role"`
}

// GetLoggedInInfoParams is just here to check if token is working
type GetLoggedInInfoParams struct {
	AuthToken string
}

// GetLoggedInInfoResult is
type GetLoggedInInfoResult struct {
	Data interface{} `json:"data"`
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
		return nil, &customError.ValidationError{
			Code:    customError.InvalidAuthToken,
			Message: err.Error(),
		}
	}
	return &GetLoggedInInfoResult{
		Data: claims,
	}, nil
}

// Login is a backend function
func (ctx *Context) Login(params LoginParams) (*LoginResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	record, err := ctx.DB.GetUserByName(params.Username)

	if err != nil {
		return nil, err
	}
	authToken, err := ctx.createToken(record.Username, record.ID, record.RoleList)
	if err != nil {
		return nil, &customError.InternalError{
			Code:    customError.UnknownError,
			Message: "Cannot generate token",
		}
	}
	return &LoginResult{
		AuthToken: authToken,
		UserID:    record.ID,
		Username:  record.Username,
		Role:      record.RoleList[0],
	}, nil
}

// Register is a backend function
func (ctx *Context) Register(params RegisterParams, role model.Role) (*RegisterResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}
	userId, err := ctx.DB.CreateUser(params.Username, role)
	if err != nil {
		return nil, &customError.UserError{
			Code:           customError.DuplicateUsername,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	return &RegisterResult{ID: userId}, nil
}
