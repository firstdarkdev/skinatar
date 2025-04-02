package main

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

var (
	// Limiter map to track rate limits for different IP addresses
	limiterMap = make(map[string]*rate.Limiter)
	limiterMux sync.Mutex
)

// Get or create a rate limiter for a given IP address
func getRateLimiter(ip string) *rate.Limiter {
	limiterMux.Lock()
	defer limiterMux.Unlock()

	// Check if we already have a limiter for this IP
	if limiter, exists := limiterMap[ip]; exists {
		return limiter
	}

	// Create a new rate limiter: 5 requests per second
	limiter := rate.NewLimiter(rate.Every(200*time.Millisecond), 1)
	limiterMap[ip] = limiter
	return limiter
}
