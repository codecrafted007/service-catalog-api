package middleware

import (
	"net/http"
	"strings"

	"github.com/codecrafted007/service-catalog-api/internal/utils"
)

func APIKeyAuth(validateKeyFunc func(string) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AuthMiddleware(w, r, next, validateKeyFunc)
		})
	}
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.Handler, validateKeyFunc func(string) bool) {
	authHeader := r.Header.Get("X-API-Key")
	if authHeader == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, nil, "API key is missing")
		return
	}

	apiKey := strings.TrimSpace(authHeader)
	if !validateKeyFunc(apiKey) {
		utils.WriteJSON(w, http.StatusForbidden, nil, "Invalid API key")
		return
	}

	next.ServeHTTP(w, r)
}
