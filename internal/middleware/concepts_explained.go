package middleware

import (
	"fmt"
	"net/http"
)

// This file explains middleware concepts with detailed examples

// Example 1: Basic Middleware Pattern
func BasicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware: Before handler")

		// This is where the magic happens - calling the next handler
		next.ServeHTTP(w, r)

		fmt.Println("Middleware: After handler")
	})
}

// Example 2: Middleware that can stop the chain
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		token := r.Header.Get("Authorization")

		if token == "" {
			// STOP the chain - don't call next.ServeHTTP
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Continue the chain
		next.ServeHTTP(w, r)
	})
}

// Example 3: Middleware that modifies the request
func AddUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token (simplified)
		// userID := "user123" // In real app, extract from JWT

		// Add user ID to request context
		// ctx := r.Context()
		// ctx = context.WithValue(ctx, "userID", userID)
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Example 4: Middleware that captures response
func ResponseCaptureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture status
		customWriter := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(customWriter, r)

		// Now we can access the status code
		fmt.Printf("Response status: %d\n", customWriter.statusCode)
	})
}

// Custom response writer to capture status codes
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}
