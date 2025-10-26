package middleware

import (
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.RWMutex
)

func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	go cleanupVisitors(window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)

			mu.Lock()
			v, exists := visitors[ip]
			now := time.Now()

			if !exists {
				visitors[ip] = &visitor{lastSeen: now, count: 1}
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if now.Sub(v.lastSeen) > window {
				v.count = 1
				v.lastSeen = now
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if v.count >= requests {
				mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			v.count++
			v.lastSeen = now
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func cleanupVisitors(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		now := time.Now()
		for ip, v := range visitors {
			if now.Sub(v.lastSeen) > interval {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
