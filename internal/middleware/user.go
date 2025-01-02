package middleware

import (
	"context"
	"net/http"
)

type key int

const UserIDKey key = 0

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user_id from the request context or token
		userID := r.Header.Get("X-User-ID") // Example: replace with actual extraction logic
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
} 