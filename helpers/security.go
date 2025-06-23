package helpers

import (
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type rateLimiter struct {
	timestamp time.Time
	count     int
}

var (
	rateLimitMap = make(map[string]*rateLimiter)
	rateMutex    sync.Mutex
)

// Check if email is valid
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

// Password must be at least 8 characters long and contain numbers and symbols.
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasNumber := strings.ContainsAny(password, os.Getenv("PASSWORD"))
	hasSymbol := strings.ContainsAny(password, os.Getenv("PASSWORD_HASH"))
	return hasNumber && hasSymbol
}

// Check the rate limit of the IP: maximum 5 requests per 10 seconds
func IsRateLimited(ip string) bool {
	rateMutex.Lock()
	defer rateMutex.Unlock()

	now := time.Now()
	limiter, exists := rateLimitMap[ip]
	if !exists || now.Sub(limiter.timestamp) > 10*time.Second {
		rateLimitMap[ip] = &rateLimiter{timestamp: now, count: 1}
		return false
	}
	if limiter.count >= 5 {
		return true
	}
	limiter.count++
	return false
}
