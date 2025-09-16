package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/model/dto"
	log "github.com/sirupsen/logrus"

	"net/http"
)

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if token == "" || token[:6] != "Bearer" {
			log.Warn("token not found or wrong format")

			response := &dto.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "token not found or wrong format",
				Result: nil,
			}

			helper.ResponseJSON(w, response)
			return
		}

		token = token[7:]

		claims, err := helper.ValidateJWTToken(token)

		if err != nil {
			log.Warn("token validation failed")

			response := &dto.WebResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
				Result: nil,
			}

			helper.ResponseJSON(w, response)
			return
		}

		memberID, err := uuid.Parse(claims.MemberID)
		if err != nil {
			log.WithError(err).Warn("MemberID is not a valid UUID")

			response := &dto.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "MemberID is not a valid UUID",
				Result: nil,
			}

			helper.ResponseJSON(w, response)
			return
		}

		role := claims.Role

		if role == "" {
			log.Warn("Role not found")

			response := &dto.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "Role not found",
				Result: nil,
			}

			helper.ResponseJSON(w, response)
			return
		}

		vals := map[string]any{
			"memberID": memberID,
			"role":     role,
		}

		ctx := context.WithValue(r.Context(), "memberDatas", vals)

		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)

	})
}
