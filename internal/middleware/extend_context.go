package middleware

import (
	"context"
	"net/http"
)

func ExtendContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithoutCancel(r.Context())

		newRq := r.WithContext(ctx)

		next.ServeHTTP(w, newRq)
	})
}
