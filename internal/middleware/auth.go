package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler, secretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Empty token", http.StatusUnauthorized)
			return 
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected method")
			}
			return ([]byte(secretKey)), nil
		})
		if err != nil {
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return 
		}
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return 
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			user_id := claims["user_id"]
			ctx := context.WithValue(r.Context(), "user_id", user_id)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}