package middleware

import (
	"log"
	"net/http"
	"strings"

	"go-api/utils"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Không có token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Token không hợp lệ 1", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		userID, err := utils.ValidateToken(token)
		if err != nil {
			log.Printf("Token không hợp lệ: %v", err)
			http.Error(w, "Token không hợp lệ 2", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(utils.WithUserID(r.Context(), userID))
		next.ServeHTTP(w, r)
	})
}
