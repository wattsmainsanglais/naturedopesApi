package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimitEntry struct {
	count      int
	windowStart time.Time
}

type RateLimiter struct {
	// Per-key tracking: 100 requests/hour
	keyLimits sync.Map

	// Per-IP tracking: 1000 requests/day
	ipLimits sync.Map

	keyLimit      int
	keyWindow     time.Duration
	ipLimit       int
	ipWindow      time.Duration
	cleanupTicker *time.Ticker
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		keyLimit:  100,
		keyWindow: time.Hour,
		ipLimit:   1000,
		ipWindow:  24 * time.Hour,
		cleanupTicker: time.NewTicker(10 * time.Minute),
	}

	// Background cleanup of old entries
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	for range rl.cleanupTicker.C {
		now := time.Now()

		// Clean up key limits
		rl.keyLimits.Range(func(key, value interface{}) bool {
			entry := value.(*rateLimitEntry)
			if now.Sub(entry.windowStart) > rl.keyWindow {
				rl.keyLimits.Delete(key)
			}
			return true
		})

		// Clean up IP limits
		rl.ipLimits.Range(func(key, value interface{}) bool {
			entry := value.(*rateLimitEntry)
			if now.Sub(entry.windowStart) > rl.ipWindow {
				rl.ipLimits.Delete(key)
			}
			return true
		})
	}
}

func (rl *RateLimiter) checkLimit(limitsMap *sync.Map, identifier string, limit int, window time.Duration) bool {
	now := time.Now()

	value, _ := limitsMap.LoadOrStore(identifier, &rateLimitEntry{
		count:      0,
		windowStart: now,
	})

	entry := value.(*rateLimitEntry)

	// Reset window if expired
	if now.Sub(entry.windowStart) >= window {
		entry.count = 0
		entry.windowStart = now
	}

	// Check if limit exceeded
	if entry.count >= limit {
		return false
	}

	// Increment count
	entry.count++
	return true
}

// RateLimitMiddleware applies rate limiting based on both API key and IP address
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract IP address (handle X-Forwarded-For for proxies)
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}

		// Check IP rate limit
		if !rl.checkLimit(&rl.ipLimits, ip, rl.ipLimit, rl.ipWindow) {
			http.Error(w, "Rate limit exceeded: Too many requests from this IP (1000/day)", http.StatusTooManyRequests)
			return
		}

		// Check API key rate limit (if key is present)
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "" {
			if !rl.checkLimit(&rl.keyLimits, apiKey, rl.keyLimit, rl.keyWindow) {
				http.Error(w, "Rate limit exceeded: Too many requests with this API key (100/hour)", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Stop stops the cleanup ticker
func (rl *RateLimiter) Stop() {
	rl.cleanupTicker.Stop()
}
