package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	rsp "zhtcloud/pkg/response"
	"zhtcloud/utils/logger"
)

// RateLimiter manages rate limiting for clients
type RateLimiter struct {
	mu              sync.RWMutex
	clients         map[string]*ClientLimiter
	maxTokens       int
	refillRate      time.Duration
	cleanupInterval time.Duration
}

// ClientLimiter tracks individual client request limits
type ClientLimiter struct {
	tokens      int
	lastRefill  time.Time
	lastRequest time.Time
}

// RateLimiterConfig holds configuration for rate limiter
type RateLimiterConfig struct {
	MaxRequests     int           // Maximum requests allowed, bucket size
	Window          time.Duration // Time window for rate limiting
	CleanupInterval time.Duration // How often to cleanup old clients
}

// NewRateLimiter creates a new rate limiter with specified config
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		clients:         make(map[string]*ClientLimiter),
		maxTokens:       config.MaxRequests,
		refillRate:      config.Window / time.Duration(config.MaxRequests),
		cleanupInterval: config.CleanupInterval,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// cleanup removes old client records to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, client := range rl.clients {
			// Remove clients that haven't made requests for 10 minutes
			if now.Sub(client.lastRequest) > 10*time.Minute {
				delete(rl.clients, key)
			}
		}
		rl.mu.Unlock()
	}
}

// isAllowed checks if a client is allowed to make a request
func (rl *RateLimiter) isAllowed(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[clientID]

	if !exists {
		// New client gets full tokens minus one for current request
		client = &ClientLimiter{
			tokens:      rl.maxTokens - 1,
			lastRefill:  now,
			lastRequest: now,
		}
		rl.clients[clientID] = client
		return true
	}

	// Update last request time
	client.lastRequest = now

	// Calculate tokens to add based on elapsed time
	elapsed := now.Sub(client.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		client.tokens = min(rl.maxTokens, client.tokens+tokensToAdd)
		client.lastRefill = now
	}

	// Check if request is allowed
	if client.tokens > 0 {
		client.tokens--
		return true
	}

	return false
}

// getClientID extracts client identifier from request
func getClientID(r *http.Request) string {
	// Priority: X-Real-IP > X-Forwarded-For > RemoteAddr
	clientIP := r.Header.Get("X-Real-IP") // TODO: IPs of hackers will change dynamically, so it may not be reliable
	if clientIP == "" {
		clientIP = r.Header.Get("X-Forwarded-For")
		if clientIP != "" {
			// X-Forwarded-For can contain multiple IPs, use the first one
			clientIP = strings.Split(clientIP, ",")[0]
			clientIP = strings.TrimSpace(clientIP)
		}
	}
	if clientIP == "" {
		clientIP = r.RemoteAddr
		// Remove port from IP:port format
		if idx := strings.LastIndex(clientIP, ":"); idx != -1 {
			clientIP = clientIP[:idx]
		}
	}

	// For specific endpoints, include additional identifiers
	if mac := r.URL.Query().Get("mac"); mac != "" {
		clientIP = clientIP + ":" + mac
	}
	if userID := r.Header.Get("User-ID"); userID != "" {
		clientIP = clientIP + ":" + userID
	}

	return clientIP
}

func RateLimiterMiddleware(config RateLimiterConfig) func(http.HandlerFunc) http.HandlerFunc {
	limiter := NewRateLimiter(config)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) { // wrapper
			clientID := getClientID(r)

			if !limiter.isAllowed(clientID) {
				rs := map[string]interface{}{
					"code": rsp.EXCEED_RATE_LIMIE,
					"msg":  rsp.CodeToMsgMap[rsp.EXCEED_RATE_LIMIE],
				}

				w.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(w).Encode(rs)
				if err != nil {
					logger.Error("JSON encode error: %v", err)
					rs["code"] = rsp.SERVER_ERROR
					rs["msg"] = rsp.CodeToMsgMap[rsp.SERVER_ERROR]
				}
				return
			}

			next(w, r) // real request handling
		}
	}
}

// Predefined configurations for common use cases
var (
	// StrictConfig for sensitive endpoints like auth
	StrictConfig = RateLimiterConfig{
		MaxRequests:     5,
		Window:          1 * time.Minute,
		CleanupInterval: 5 * time.Minute,
	}

	// ModerateConfig for regular API endpoints
	ModerateConfig = RateLimiterConfig{
		MaxRequests:     30,
		Window:          1 * time.Minute,
		CleanupInterval: 10 * time.Minute,
	}

	// LenientConfig for high-frequency endpoints
	LenientConfig = RateLimiterConfig{
		MaxRequests:     100,
		Window:          1 * time.Minute,
		CleanupInterval: 15 * time.Minute,
	}
)

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
