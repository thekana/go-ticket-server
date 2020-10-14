package response

import "net/http"

type statusCodeRecorder struct {
	http.ResponseWriter
	// http.Hijacker
	StatusCode int
}
