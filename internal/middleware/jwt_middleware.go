package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

func MiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("=== JWT Middleware called for: %s ===", r.URL.Path)

		// Get secret key from environment variable
		secret := os.Getenv("GYM_APP_SECRET_KEY")
		if secret == "" {
			log.Fatal("GYM_APP_SECRET_KEY environment variable not set")
		}

		//logic for jwt extraction from authorization header and parsing and verification
		authHeader := r.Header.Get("Authorization")
		// log.Printf("Authorization header: %s", authHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Missing or invalid authorization header")
			http.Error(w, "Missing or invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Token string: %s", tokenString[:20]+"...") // Log first 20 chars for security

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "Unexpected signing method", http.StatusNotAcceptable)
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil // Use environment variable instead of hardcoded string
		})

		if err != nil {
			log.Printf("JWT parsing error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Printf("Token is not valid")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Token is valid, extracting claims...")

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Printf("Claims extracted successfully: %+v", claims)
			ctx := context.WithValue(r.Context(), ContextKey("user"), claims)
			r = r.WithContext(ctx)
			log.Printf("Context set successfully")
		} else {
			log.Printf("Failed to extract claims from token")
		}

		next.ServeHTTP(w, r)
	})
}

func GenerateToken(username string) (string, error) {
	secret := os.Getenv("GYM_APP_SECRET_KEY")
	if secret == "" {
		return "", fmt.Errorf("invalid_secret")
	}
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "the-gym-app",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//signing the token with the secret
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	//returning signed token
	return signedToken, nil
}
