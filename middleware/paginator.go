package middleware

import (
	"context"
	"net/http"

	"github.com/nanoLeinz/librarium/helper"
)

func Paginator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		page := r.URL.Query().Get("page")
		pageSize := r.URL.Query().Get("page_size")

		ctx := context.WithValue(r.Context(), helper.KeyCon("page"), page)
		ctx = context.WithValue(ctx, helper.KeyCon("page_size"), pageSize)
		newRq := r.WithContext(ctx)

		next.ServeHTTP(w, newRq)
	})
}
