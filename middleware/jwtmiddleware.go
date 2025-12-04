package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddlware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//get token from authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"message":"missing authorization"}`, http.StatusUnauthorized)
			return
		}

		//checkk bearear token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"message":"invalid authorization header"}`, http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		secret := os.Getenv("JWT_TOKEN")
		if secret == "" {
			secret = "default_secret"
		}

		//validate the jwt
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("umexpected sigining method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"message":"invlaid token"}`, http.StatusUnauthorized)
			return
		}
		//vlaid then procedd to the next handler
		next(w, r)
	}
}
