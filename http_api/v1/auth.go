package v1

import (
	"net/http"

	"ticket-reservation/app"
	"ticket-reservation/http_api/routes"
)

var AuthRoutes = routes.Routes{
	routes.Route{
		Name:        "Login",
		Path:        "/login",
		Method:      "POST",
		HandlerFunc: Login,
	},
}

func init() {
	RouteDefinitions = append(RouteDefinitions, routes.RouteDefinition{
		Routes: AuthRoutes,
		Prefix: "",
	})
}

// FIXME
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
		return err
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

	return nil
}
