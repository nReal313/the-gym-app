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

type contextKey string

func MiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get secret key from environment variable
		secret := os.Getenv("GYM_APP_SECRET_KEY")
		if secret == "" {
			log.Fatal("GYM_APP_SECRET_KEY environment variable not set")
		}

		//logic for jwt extraction from authorization header and parsing and verification
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "Unexpected signing method", http.StatusNotAcceptable)
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil // Use environment variable instead of hardcoded string
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), contextKey("user"), claims)
			r = r.WithContext(ctx)
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
