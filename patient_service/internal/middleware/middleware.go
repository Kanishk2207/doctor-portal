package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"patient_service/internal/configs"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

var jwtKey = []byte(configs.LoadConfig().JWTSecret)

func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/signup" || r.URL.Path == "/login" || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Missing Authorization header")
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Invalid token format")
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			var validationErr *jwt.ValidationError
			if errors.As(err, &validationErr) {
				if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
					log.Printf("Token expired: %v", err)
					http.Error(w, "Unauthorized: Token expired", http.StatusUnauthorized)
					return
				}
				log.Printf("Token validation error: %v", err)
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "Unauthorized: Token error", http.StatusUnauthorized)
			return
		}

		if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
			log.Printf("Token has expired")
			http.Error(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
			return
		}

		delete(claims, "exp")

		ctx := context.WithValue(r.Context(), contextKey("userPayload"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ExtractClaimsFromContext(r *http.Request) (jwt.MapClaims, error) {
	claims, ok := r.Context().Value(contextKey("userPayload")).(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("no valid token payload in context")
	}
	return claims, nil
}

func AttachClaimsToContext(r *http.Request, claims jwt.MapClaims) *http.Request {
	ctx := context.WithValue(r.Context(), contextKey("userPayload"), claims)
	return r.WithContext(ctx)
}
