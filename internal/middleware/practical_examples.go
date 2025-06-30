package middleware

import (
	"net/http"
	"time"
)

// CORS Middleware - handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Rate Limiting Middleware (simplified)
func RateLimitMiddleware(next http.Handler) http.Handler {
	// In a real app, you'd use a proper rate limiter
	// This is just a conceptual example
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user has exceeded rate limit
		// For now, we'll just add a small delay to simulate processing
		time.Sleep(10 * time.Millisecond)

		next.ServeHTTP(w, r)
	})
}

// Security Headers Middleware
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}

// Request ID Middleware - adds unique ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a simple request ID (in real app, use UUID)
		requestID := time.Now().Format("20060102150405")

		// Add to response headers
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}
