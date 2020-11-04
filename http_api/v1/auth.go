package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"ticket-reservation/app"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
	"ticket-reservation/http_api/response"
	"ticket-reservation/http_api/routes"
)

// AuthRoutes is for adding auth api routes
var AuthRoutes = routes.Routes{
	routes.Route{
		Name:        "Login",
		Path:        "/login",
		Method:      "POST",
		HandlerFunc: Login,
	},
	routes.Route{
		Name:        "Customer Register",
		Path:        "/register",
		Method:      "POST",
		HandlerFunc: Register,
	},
	routes.Route{
		Name:        "Check",
		Path:        "/check",
		Method:      "GET",
		HandlerFunc: GetLoggedInInfo,
	},
}

func init() {
	RouteDefinitions = append(RouteDefinitions, routes.RouteDefinition{
		Routes: AuthRoutes,
		Prefix: "",
	})
}

func GetLoggedInInfo(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.GetLoggedInInfoParams
	var err error
	input.AuthToken, err = extractBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		return err
	}
	resData, err := ctx.GetLoggedInInfo(input)
	if err != nil {
		return err
	}

	data, err := json.Marshal(&response.Response{
		Code:    0,
		Message: "",
		Data:    resData.Data,
	})
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	return err
}

func Login(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.LoginParams

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return &customError.UserError{
			Code:           customError.InvalidJSONString,
			Message:        "Invalid JSON string",
			HTTPStatusCode: http.StatusBadRequest,
		}
	}

	resData, err := ctx.Login(input)
	if err != nil {
		return &customError.UserError{
			Code:           customError.UserNotFound,
			Message:        fmt.Sprint(err),
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	data, err := json.Marshal(&response.Response{
		Code:    0,
		Message: "",
		Data:    resData,
	})
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	return err
}

func Register(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.RegisterParams
	var role = model.Customer
	if roleType := r.URL.Query().Get("type"); strings.ToLower(roleType) == "organizer" {
		role = model.Organizer
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return &customError.UserError{
			Code:           customError.InvalidJSONString,
			Message:        "Invalid JSON string",
			HTTPStatusCode: http.StatusBadRequest,
		}
	}

	resData, err := ctx.Register(input, role)
	if err != nil {
		return &customError.UserError{
			Code:           customError.DuplicateUsername,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	data, err := json.Marshal(&response.Response{
		Code:    0,
		Message: "",
		Data:    resData,
	})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	w.WriteHeader(http.StatusCreated)
	return err
}
