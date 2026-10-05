package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Visitor struct {
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     float64
	burst    float64
}

func NewRateLimiter(requestsPerMinute int, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     float64(requestsPerMinute) / 60.0,
		burst:    float64(burst),
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &Visitor{
			tokens:     rl.burst,
			lastRefill: time.Now(),
		}
		rl.visitors[ip] = v
	}
	return v
}

func (rl *RateLimiter) allow(ip string) bool {
	v := rl.getVisitor(ip)
	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(v.lastRefill).Seconds()

	v.tokens += elapsed * rl.rate
	if v.tokens > rl.burst {
		v.tokens = rl.burst
	}
	v.lastRefill = now

	if v.tokens >= 1 {
		v.tokens--
		return true
	}
	return false
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			v.mu.Lock()
			inactive := time.Since(v.lastRefill) > 10*time.Minute
			v.mu.Unlock()
			if inactive {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		if !rl.allow(ip) {
			w.Header().Set("Retry-After", "60")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"Too many requests, slow down."}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
