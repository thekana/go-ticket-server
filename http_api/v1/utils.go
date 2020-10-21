package v1

import (
	"fmt"
	"net/http"
	"ticket-reservation/app"
	"ticket-reservation/http_api/routes"
)

// AuthRoutes is for adding auth api routes
var UtilRoutes = routes.Routes{
	routes.Route{
		Name:        "Print",
		Path:        "/print",
		Method:      "GET",
		HandlerFunc: PrintSystem,
	},
	routes.Route{
		Name:        "Populate",
		Path:        "/pop",
		Method:      "GET",
		HandlerFunc: Populate,
	},
}

func init() {
	RouteDefinitions = append(RouteDefinitions, routes.RouteDefinition{
		Routes: UtilRoutes,
		Prefix: "",
	})
}

func PrintSystem(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Print(`
  ______  __      __   ______  ________  ________  __       __ 
 /      \|  \    /  \ /      \|        \|        \|  \     /  \
|  $$$$$$\\$$\  /  $$|  $$$$$$\\$$$$$$$$| $$$$$$$$| $$\   /  $$
| $$___\$$ \$$\/  $$ | $$___\$$  | $$   | $$__    | $$$\ /  $$$
 \$$    \   \$$  $$   \$$    \   | $$   | $$  \   | $$$$\  $$$$
 _\$$$$$$\   \$$$$    _\$$$$$$\  | $$   | $$$$$   | $$\$$ $$ $$
|  \__| $$   | $$    |  \__| $$  | $$   | $$_____ | $$ \$$$| $$
 \$$    $$   | $$     \$$    $$  | $$   | $$     \| $$  \$ | $$
  \$$$$$$     \$$      \$$$$$$    \$$    \$$$$$$$$ \$$      \$$
`)
	ctx.DB.PrintSystem()
	w.WriteHeader(http.StatusOK)
	return nil
}

var did bool

func Populate(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	if did {
		w.WriteHeader(http.StatusNotModified)
	} else {
		ctx.DB.PopulateSystem()
		did = true
		w.WriteHeader(http.StatusOK)
	}
	return nil
}
