package middleware

import (
	"cometosee/common"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastseen time.Time
}

var (
	mu        sync.Mutex
	clients   = make(map[string]*client)
	ratelimit = rate.Every(10 * time.Hour)
	burst     = 5000
)

func getClients(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if c, exists := clients[ip]; exists {
		c.lastseen = time.Now()
		return c.limiter
	}

	limiter := rate.NewLimiter(ratelimit, burst)
	clients[ip] = &client{limiter: limiter, lastseen: time.Now()}
	return limiter
}

func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastseen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	go cleanupClients()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getClients(ip)

		if !limiter.Allow() {
			common.WriteJSONError(w, "too many request", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
