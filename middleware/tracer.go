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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("generating traceID")

		const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		ran := rand.New(rand.NewSource(time.Now().UnixNano()))
		trace := make([]byte, 6)
		for i := range trace {
			trace[i] = charset[ran.Intn(len(charset))]
		}

		traceID := string(trace)

		log.WithField("traceID", string(traceID)).Info("traceID generated")
		log.WithField("traceID", traceID).Info("adding traceID to contect")

		ctx := context.WithValue(r.Context(), helper.KeyCon("traceID"), traceID)
		newR := r.WithContext(ctx)

		next.ServeHTTP(w, newR)
	})
}
