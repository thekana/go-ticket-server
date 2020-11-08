package http_api

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"

	"ticket-reservation/app"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/http_api/middleware"
	"ticket-reservation/http_api/response"
	"ticket-reservation/http_api/routes"
	apiv1 "ticket-reservation/http_api/v1"
	"ticket-reservation/log"
	"ticket-reservation/prometheus"
	"ticket-reservation/utils"
)

var routeDefinitions = make([]routes.RouteDefinition, 0)

type statusCodeRecorder struct {
	http.ResponseWriter
	// http.Hijacker
	StatusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

type API struct {
	App    *app.App
	Config *Config
}

func New(app *app.App) (api *API, err error) {
	api = &API{App: app}
	api.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}
	return api, nil
}

func readAccessTokenSignerPublicKey(filepath string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	publicKey := utils.BytesToPublicKey(publicKeyBytes)
	return publicKey, nil
}

func (api *API) Init() http.Handler {
	router := mux.NewRouter()

	router.Use(prometheusMiddleware)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	apiRouter := router.PathPrefix("/api").Subrouter()

	for _, routeDefinition := range routeDefinitions {
		api.addRoutes(apiRouter, routeDefinition.Routes, routeDefinition.Prefix)
	}

	// API v1
	subRouter := apiRouter.PathPrefix("/v1").Subrouter()
	for _, routeDefinition := range apiv1.RouteDefinitions {
		api.addRoutes(subRouter, routeDefinition.Routes, routeDefinition.Prefix)
	}

	return middleware.RequestIDMiddleware(api.loggingMiddleware(router))
}

// addRoutes is an internal function to add routes to the api app.
func (api *API) addRoutes(router *mux.Router, routes routes.Routes, prefix string) {
	if prefix == "" {
		for _, route := range routes {
			router.
				Handle(route.Path, api.handler(route.HandlerFunc)).
				Methods(route.Method).
				Name(route.Name)
		}
	} else {
		subRouter := router.PathPrefix(prefix).Subrouter()
		for _, route := range routes {
			subRouter.
				Handle(route.Path, api.handler(route.HandlerFunc)).
				Methods(route.Method).
				Name(route.Name)
		}
	}
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		startTime := time.Now()

		next.ServeHTTP(w, r)
		duration := time.Since(startTime)

		statusCode := w.(*statusCodeRecorder).StatusCode

		prometheus.RecordHTTPAPICallDurationMetrics(duration, statusCode, path)
	})
}

func (api *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		ctx := r.Context()

		requestID := middleware.GetRequestID(ctx)

		logger := api.App.Logger.WithFields(log.Fields{
			"module":     "api",
			"remote":     api.RemoteAddressForRequest(r),
			"request_id": requestID,
		})

		w = &statusCodeRecorder{
			ResponseWriter: w,
		}

		if api.Config.LogLevel == log.Debug {
			// Save a copy of this request for debugging.
			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				logger.Errorf("%+v", err)
			}
			logger.Debugf("HTTP Request:\n%s", string(requestDump))

			rec := httptest.NewRecorder()
			next.ServeHTTP(&statusCodeRecorder{
				ResponseWriter: rec,
			}, r)

			responseDump, err := httputil.DumpResponse(rec.Result(), true)
			if err != nil {
				logger.Errorf("%+v", err)
			}
			logger.Debugf("HTTP Response:\n%s", string(responseDump))

			// copy the captured response headers to new response
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			_, _ = rec.Body.WriteTo(w)
		} else {
			next.ServeHTTP(w, r)
		}

		duration := time.Since(startTime)
		statusCode := w.(*statusCodeRecorder).StatusCode
		logger = logger.WithFields(log.Fields{
			"duration":    duration.String(),
			"status_code": statusCode,
		})
		logger.Infof("%s %s", r.Method, r.URL.RequestURI())
	})
}

func (api *API) handler(f func(*app.Context, http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := middleware.GetRequestID(ctx)

		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)

		appCtx := api.App.NewContext().WithRemoteAddress(api.RemoteAddressForRequest(r))

		appCtx = appCtx.WithLogger(appCtx.Logger.WithFields(
			log.Fields{
				"request_id": requestID,
			},
		))

		logger := appCtx.Logger.WithFields(log.Fields{
			"module": "api",
			"remote": appCtx.RemoteAddress,
		})

		w.Header().Set("Content-Type", "application/json")

		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("%v: %s", r, debug.Stack())
				resData, err := json.Marshal(&response.Response{
					Code:    customError.UnknownError,
					Message: "Internal server error",
				})
				if err != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write(resData)
			}
		}()

		if err := f(appCtx, w, r); err != nil {
			if verr, ok := err.(*customError.ValidationError); ok {
				data, err := json.Marshal(&response.Response{
					Code:    verr.Code,
					Message: verr.Message,
				})
				if err == nil {
					w.WriteHeader(http.StatusBadRequest)
					_, err = w.Write(data)
				}
				if err != nil {
					logger.Errorf("HTTP hanlder validation err: %+v", err)
				}
			} else if uerr, ok := err.(*customError.UserError); ok {
				if uerr.Code == 0 {
					uerr.Code = customError.UnknownError
				}
				data, err := json.Marshal(&response.Response{
					Code:    uerr.Code,
					Message: uerr.Message,
				})
				if err == nil {
					if uerr.HTTPStatusCode == 0 {
						w.WriteHeader(http.StatusBadRequest)
					} else {
						w.WriteHeader(uerr.HTTPStatusCode)
					}
					_, err = w.Write(data)
				}
				if err != nil {
					logger.Errorf("HTTP handler user err: %+v", err)
				}
			} else if aerr, ok := err.(*customError.AuthorizationError); ok {
				data, err := json.Marshal(&response.Response{
					Code:    aerr.Code,
					Message: aerr.Message,
				})
				if err == nil {
					w.WriteHeader(http.StatusUnauthorized)
					_, err = w.Write(data)
				}
				if err != nil {
					logger.Errorf("%+v", err)
				}
			} else if ierr, ok := err.(*customError.InternalError); ok {
				data, err := json.Marshal(&response.Response{
					Code:    ierr.Code,
					Message: ierr.Message,
				})
				if err == nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, err = w.Write(data)
				}
				if err != nil {
					logger.Errorf("%+v", err)
				}
			} else {
				logger.Errorf("HTTP handler err: %+v", err)
			}
		}
		statusCode := w.(*statusCodeRecorder).StatusCode
		if statusCode == 0 {
			logger.Errorf("return HTTP status code not set, responding with 500 internal server error")
			resData, err := json.Marshal(&response.Response{
				Code:    customError.UnknownError,
				Message: "Internal server error",
			})
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write(resData)
		}
	})
}

func (api *API) RemoteAddressForRequest(r *http.Request) string {
	return r.RemoteAddr
}
