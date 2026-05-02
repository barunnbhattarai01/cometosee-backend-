package middleware

import (
	"context"
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
		secret := os.Getenv("SECRET")

		//validate the jwt
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("umexpected sigining method")
			}
			return []byte(secret), nil

		})

		if err != nil || !token.Valid {
			http.Error(w, `{"message":"invlaid token"}`, http.StatusUnauthorized)
			fmt.Println("token error:", err)
			return
		}

		//  extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"message":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		//  get email from claims
		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, `{"message":"email not found in token"}`, http.StatusUnauthorized)
			return
		}

		// add email to context
		ctx := context.WithValue(r.Context(), "email", email)

		next(w, r.WithContext(ctx))

	}
}
