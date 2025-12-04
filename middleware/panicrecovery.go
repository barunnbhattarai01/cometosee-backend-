package middleware

import (
	"cometosee/common"
	"log"
	"net/http"
)

func PanicrecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("recovery from error :%v", err)
				common.WriteJSONError(w, "Internal server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
