package middleware

import (
	"net/http"
	"slices"

	"helpdesk-api/internal/model"
)

func RequireRole(roles ...model.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*Claims)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !slices.Contains(roles, model.Role(claims.Role)) {
				w.WriteHeader(http.StatusForbidden) // 403
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
