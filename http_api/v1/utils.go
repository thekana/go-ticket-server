package v1

import (
	"net/http"
	"ticket-reservation/app"
	"ticket-reservation/http_api/routes"
)

var UtilRoutes = routes.Routes{
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

func Populate(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	err := ctx.RedisCache.SetNXEventQuota(-1, -1)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
	} else {
		ctx.DB.PopulateSystem()
		w.WriteHeader(http.StatusOK)
	}
	return nil
}
