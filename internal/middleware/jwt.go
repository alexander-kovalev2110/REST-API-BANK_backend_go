package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const CustomerIDKey contextKey = "customer_id"

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Authorization header missing"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"Invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"Invalid token claims"}`, http.StatusUnauthorized)
				return
			}

			rawID, ok := claims["customer_id"]
			if !ok {
				http.Error(w, `{"error":"Customer ID not found in token"}`, http.StatusUnauthorized)
				return
			}

			var customerID int64
			switch v := rawID.(type) {
			case float64:
				customerID = int64(v)
			case int64:
				customerID = v
			default:
				http.Error(w, `{"error":"Invalid customer_id format"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CustomerIDKey, customerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetCustomerID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(CustomerIDKey).(int64)
	return id, ok
}
