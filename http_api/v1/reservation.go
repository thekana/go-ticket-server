package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"ticket-reservation/app"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/http_api/response"
	"ticket-reservation/http_api/routes"
)

var ReservationRoutes = routes.Routes{
	routes.Route{
		Name:        "Make a reservation",
		Path:        "/reserve",
		Method:      "POST",
		HandlerFunc: Reserve,
	},
	routes.Route{
		Name:        "Users view their reservations",
		Path:        "/view",
		Method:      "GET",
		HandlerFunc: ViewAllReservations,
	},
	routes.Route{
		Name:        "Cancel a reservation",
		Path:        "/cancel",
		Method:      "POST",
		HandlerFunc: CancelReservation,
	},
}

func init() {
	RouteDefinitions = append(RouteDefinitions, routes.RouteDefinition{
		Routes: ReservationRoutes,
		Prefix: "/reservation",
	})
}

func Reserve(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.MakeReservationParams
	var err error
	input.AuthToken, err = extractBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		return err
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
	resData, err := ctx.MakeReservation(input)
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

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	return err
}

func ViewAllReservations(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.ViewReservationsParams
	var err error
	input.AuthToken, err = extractBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		return err
	}

	resData, err := ctx.ViewReservations(input)
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
}

func CancelReservation(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.CancelReservationParams
	var err error
	input.AuthToken, err = extractBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		return err
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
	resData, err := ctx.CancelReservation(input)
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
}
