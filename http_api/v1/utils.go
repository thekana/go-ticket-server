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

var did bool

func Populate(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	if did {
		w.WriteHeader(http.StatusNotModified)
	} else {
		ctx.DB.PopulateSystem()
		//ctx.My.EventQuotaMap.Set(1, 10000)
		//ctx.My.EventQuotaMap.Set(2, 10000)
		//ctx.My.EventQuotaMap.Set(3, 10000)
		//ctx.My.EventQuotaMap.Set(4, 10000)
		did = true
		w.WriteHeader(http.StatusOK)
	}
	return nil
}
