package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a package-local type for context keys to avoid collisions.
type contextKey string

// JWTAuth validates a Supabase JWT and extracts org_id into context.
// Register this instead of OrgID middleware once Supabase Auth is fully wired.
// For now OrgID (reading X-Org-ID header) is used in development.
func JWTAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"invalid claims"}`, http.StatusUnauthorized)
				return
			}

			orgID, _ := claims["org_id"].(string)
			if orgID == "" {
				http.Error(w, `{"error":"missing org_id claim"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), contextKey("org_id"), orgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
