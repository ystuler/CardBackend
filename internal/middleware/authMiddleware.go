package middleware

import (
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
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			http.Error(w, "token is empty", http.StatusUnauthorized)
			return
		}

		claims, err := util.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			http.Error(w, "invalid user id in token", http.StatusUnauthorized)
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
