package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/nanoLeinz/librarium/helper"
	log "github.com/sirupsen/logrus"
)

func GenerateTraceID(next http.Handler) http.Handler {
	log.Info("generating traceID")

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	traceID := make([]byte, 6)
	for i := range traceID {
		traceID[i] = charset[r.Intn(len(charset))]
	}

	log.WithField("traceID", traceID).Info("traceID generated")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("traceID", traceID).Info("adding traceID to contect")

		ctx := context.WithValue(r.Context(), helper.KeyCon("traceID"), traceID)
		newR := r.WithContext(ctx)

		next.ServeHTTP(w, newR)
	})
}
