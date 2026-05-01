package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func JWTProtected(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read token from cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				JSON(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
				return
			}

			// parse and validate token
			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
			}

			role, ok := claims["role"].(string)
			if !ok {
				JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			// parse user id into pgtype.UUID
			var userID pgtype.UUID
			if err := userID.Scan(claims["user_id"]); err != nil {
				JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				return
			}

			// store in context
			ctx := context.WithValue(r.Context(), ContextUserID, userID)
			ctx = context.WithValue(ctx, ContextRole, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(ContextRole) != "admin" {
			JSON(w, http.StatusForbidden, map[string]string{"error": "admin only"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SupplierOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(ContextRole)
		if role != "supplier" && role != "admin" {
			JSON(w, http.StatusForbidden, map[string]string{"error": "supplier only"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
