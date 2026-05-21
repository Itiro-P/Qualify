package middleware

import (
	"fmt"
	"main/internal/utils"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter tracks request counts per IP
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*ClientRateLimit
	limit   int
	window  time.Duration
}

// ClientRateLimit tracks rate limit info for a client
type ClientRateLimit struct {
	count       int
	lastReset   time.Time
	lastAttempt time.Time
	locked      bool
	lockTime    time.Time
}

// NewRateLimiter creates a new rate limiter
// limit: max requests per window
// window: time duration for the limit window

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ClientRateLimit),
		limit:   limit,
		window:  window,
	}
	// limpa entradas antigas a cada 5 minutos
	go func() {
		for range time.Tick(5 * time.Minute) {
			rl.mu.Lock()
			for ip, c := range rl.clients {
				if time.Since(c.lastAttempt) > rl.window*2 {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

// IsAllowed checks if a request is allowed
func (rl *RateLimiter) IsAllowed(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	client, exists := rl.clients[clientIP]

	now := time.Now()

	if !exists {
		rl.clients[clientIP] = &ClientRateLimit{
			count:       1,
			lastReset:   now,
			lastAttempt: now,
		}
		return true
	}

	// Check if locked
	if client.locked {
		if now.Sub(client.lockTime) > 15*time.Minute {
			client.locked = false
			client.count = 0
		} else {
			return false
		}
	}

	// Check if window expired
	if now.Sub(client.lastReset) > rl.window {
		client.count = 1
		client.lastReset = now
		client.lastAttempt = now
		return true
	}

	client.lastAttempt = now

	// Increment and check limit
	client.count++
	if client.count > rl.limit {
		client.locked = true
		client.lockTime = now
		return false
	}

	return true
}

// RateLimitMiddleware creates a gin middleware for rate limiting
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := getClientIP(c)

		if !rl.IsAllowed(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "too many requests",
				"message": "Rate limit exceeded. Please try again later.",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS (Strict-Transport-Security)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:;")

		c.Next()
	}
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	jwtManager := utils.NewJWTManager()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "missing authorization header",
				"message": "Authorization header is required",
				"code":    "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid authorization header",
				"message": "Authorization header must be in format: Bearer <token>",
				"code":    "INVALID_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid token",
				"message": err.Error(),
				"code":    "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Next()
	}
}

// ErrorHandlingMiddleware handles panics and errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":     "internal server error",
					"message":   "An unexpected error occurred",
					"code":      "INTERNAL_ERROR",
					"timestamp": time.Now(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// LoggingMiddleware logs request information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.RequestURI
		clientIP := getClientIP(c)

		// Log request details
		fmt.Printf("[%s] %s %s - Status: %d - Duration: %v - IP: %s\n",
			time.Now().Format(time.RFC3339),
			method,
			path,
			statusCode,
			duration,
			clientIP,
		)
	}
}

// Helper function to get client IP
func getClientIP(c *gin.Context) string {
	// Só use X-Forwarded-For se sua infra garante que o proxy seta esse header.
	/**
	if forwarded := c.Request.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Take first IP if multiple are present
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if realIP := c.Request.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	**/
	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}
