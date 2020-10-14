package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

type Id []byte

func NewId() Id {
	ret := make(Id, 20)
	if _, err := rand.Read(ret); err != nil {
		panic(err)
	}
	return ret
}

type contextKey string

const contextKeyRequestID contextKey = "requestID"

func assignRequestID(ctx context.Context) context.Context {
	requestID := base64.RawURLEncoding.EncodeToString(NewId())

	return context.WithValue(ctx, contextKeyRequestID, requestID)
}

func GetRequestID(ctx context.Context) string {
	requestID := ctx.Value(contextKeyRequestID)

	if ret, ok := requestID.(string); ok {
		return ret
	}

	return ""
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r = r.WithContext(assignRequestID(ctx))

		next.ServeHTTP(w, r)
	})
}
