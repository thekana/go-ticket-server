package http_api

import (
	"net/http"

	"ticket-reservation/app"
	"ticket-reservation/http_api/routes"
)

var HealthRoutes = routes.Routes{
	routes.Route{
		Name:        "Health",
		Path:        "/health",
		Method:      "GET",
		HandlerFunc: Health,
	},
}

func init() {
	routeDefinitions = append(routeDefinitions, routes.RouteDefinition{
		Routes: HealthRoutes,
		Prefix: "",
	})
}

func Health(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(""))
	return err
}
