package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware is a simple middleware that logs request details
// This demonstrates the basic middleware pattern: func(http.Handler) http.Handler
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do something BEFORE the next handler
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Do something AFTER the next handler
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}
