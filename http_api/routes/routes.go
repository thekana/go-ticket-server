package routes

import (
	"net/http"

	"ticket-reservation/app"
)

// Route is used to define http routes for the app
type Route struct {
	Name        string
	Path        string
	Method      string
	HandlerFunc func(*app.Context, http.ResponseWriter, *http.Request) error
}

// Routes is a collection of multiple http Routes
type Routes []Route

// RouteDefinition is something
type RouteDefinition struct {
	Routes Routes
	Prefix string
}
