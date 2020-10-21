package app

import (
	"fmt"
	"strings"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type Auth struct {
	User        *model.User
	IsAdmin     bool
	IsOrganizer bool
	IsCustomer  bool
}

func roleNumToString(role model.Role) string {
	switch role {
	case model.Admin:
		return "admin"
	case model.Organizer:
		return "organizer"
	case model.Customer:
		return "customer"
	}
	return ""
}

func (ctx *Context) authorizeUser(authToken string, allowedRoles []model.Role) (*Auth, error) {
	logger := ctx.getLogger()

	tokenValid, jwtClaims, err := ctx.verifyToken(authToken)
	if err != nil {
		if err == ErrTokenExpired {
			return nil, &customError.AuthorizationError{
				Code:    customError.AuthTokenExpired,
				Message: "token expired",
			}
		}
		return nil, err
	}

	logger.Debugf("token valid: %t", tokenValid)

	if !tokenValid {
		return nil, &customError.AuthorizationError{
			Code:    customError.InvalidAuthToken,
			Message: "invalid token",
		}
	}
	roleString := fmt.Sprint((*jwtClaims)["role"])

	var permit bool
	for _, allowedRole := range allowedRoles {
		permit = strings.Contains(roleString, roleNumToString(allowedRole))
		if permit {
			break
		}
	}
	if !permit {
		return nil, &customError.AuthorizationError{
			Code:    customError.Unauthorized,
			Message: "unauthorized",
		}
	}
	auth := &Auth{
		User: &model.User{
			ID:       int64((*jwtClaims)["uid"].(float64)),
			Username: (*jwtClaims)["name"].(string),
		},
		IsAdmin:     strings.Contains(roleString, roleNumToString(model.Admin)),
		IsOrganizer: strings.Contains(roleString, roleNumToString(model.Organizer)),
		IsCustomer:  strings.Contains(roleString, roleNumToString(model.Customer)),
	}
	return auth, nil
}
