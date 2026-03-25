package middleware

import (
	"back/internal/exceptions"
	"back/internal/util"
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = 0
)

func UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrUnauthorized)
			return
		}

		claims, err := util.ParseToken(headerParts[1])
		if err != nil {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
			return
		}

		userId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserId(ctx context.Context) (int, error) {
	id, ok := ctx.Value(userCtx).(int)
	if !ok {
		return 0, errors.New("user id not found")
	}

	return id, nil
}
