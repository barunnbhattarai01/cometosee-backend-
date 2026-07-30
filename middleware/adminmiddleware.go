package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AdminJwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		header := r.Header.Get("Authorization")
		if header == "" {
			if token := r.URL.Query().Get("token"); token != "" {
				header = "Bearer " + token
			}
		}

		parts := strings.Fields(header)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"message":"invalid authorization header"}`, http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("ADMIN_SECRET")
		if secret == "" {
			secret = os.Getenv("SECRET")
		}
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, `{"message":"invalid admin token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		role, roleOK := claims["role"].(string)
		if !ok || !roleOK || role != "admin" {
			http.Error(w, `{"message":"admin access required"}`, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "adminId", claims["adminId"])
		ctx = context.WithValue(ctx, "adminEmail", claims["email"])
		next(w, r.WithContext(ctx))
	}
}
