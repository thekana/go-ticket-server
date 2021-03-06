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

// EventRoutes is for adding auth api routes
var EventRoutes = routes.Routes{
	routes.Route{
		Name:        "Create event",
		Path:        "/create",
		Method:      "POST",
		HandlerFunc: CreateEvent,
	},
	routes.Route{
		Name:        "View all event",
		Path:        "/all",
		Method:      "GET",
		HandlerFunc: ViewAllEvents,
	},
	routes.Route{
		Name:        "View a particular event",
		Path:        "/view",
		Method:      "POST",
		HandlerFunc: ViewOneEvent,
	},
	routes.Route{
		Name:        "Edit an event",
		Path:        "/edit",
		Method:      "POST",
		HandlerFunc: EditEvent,
	},
	routes.Route{
		Name:        "Delete an event",
		Path:        "/delete",
		Method:      "POST",
		HandlerFunc: DeleteEvent,
	},
}

func init() {
	RouteDefinitions = append(RouteDefinitions, routes.RouteDefinition{
		Routes: EventRoutes,
		Prefix: "/events",
	})
}

func CreateEvent(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.CreateEventParams
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
	resData, err := ctx.CreateEvent(input)
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

func ViewAllEvents(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.ViewAllEventParams
	var err error
	input.AuthToken, err = extractBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		return err
	}

	resData, err := ctx.GetAllEventDetails(input)
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

func ViewOneEvent(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.ViewEventParams
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
	resData, err := ctx.GetEventDetail(input)
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

func EditEvent(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.EditEventParams
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
	resData, err := ctx.EditEventDetail(input)
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

func DeleteEvent(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input app.DeleteEventParams
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
	resData, err := ctx.DeleteEvent(input)
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
