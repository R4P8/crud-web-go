package utils

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	timestamp time.Time
	count     int
}

var (
	rateLimitMap = make(map[string]*RateLimiter)
	rateMutex    sync.Mutex
)

// IsRateLimited checks if the IP has exceeded the rate limit (5 requests per 10 seconds)
func IsRateLimited(ip string) bool {
	rateMutex.Lock()
	defer rateMutex.Unlock()

	const (
		limitDuration = 10 * time.Second
		maxRequests   = 5
	)

	limiter, exists := rateLimitMap[ip]
	now := time.Now()

	if !exists || now.Sub(limiter.timestamp) > limitDuration {
		rateLimitMap[ip] = &RateLimiter{timestamp: now, count: 1}
		return false
	}

	if limiter.count >= maxRequests {
		return true
	}

	limiter.count++
	return false
}

// GetClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
