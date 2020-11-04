package v1

import (
	"net/http"
	"ticket-reservation/custom_error"
	"ticket-reservation/http_api/routes"
)

var RouteDefinitions = make([]routes.RouteDefinition, 0)

func extractBearerToken(bearer string) (string, error) {
	if len(bearer) < 8 {
		return "", &custom_error.AuthorizationError{
			Code:           0,
			Message:        "Invalid Token",
			HTTPStatusCode: http.StatusUnauthorized,
		}
	}
	return bearer[7:], nil
}
