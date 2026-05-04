package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
)

func WriteJSON[T any](w http.ResponseWriter, status int, resp response.APIResponse[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, response.Failure[string](msg))
}

func GetUserID(ctx context.Context) (pgtype.UUID, bool) {
	id, ok := ctx.Value(ContextUserID).(pgtype.UUID)
	return id, ok
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextRole).(string)
	return role, ok
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
			// Get token from cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "missing token")
				return
			}

			// Parse token
			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
				// Ensure correct signing method
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenInvalidClaims
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				WriteError(w, http.StatusUnauthorized, "invalid claims")
				return
			}

			// Check expiry
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					WriteError(w, http.StatusUnauthorized, "token expired")
					return
				}
			}

			// Extract role
			role, ok := claims["role"].(string)
			if !ok {
				WriteError(w, http.StatusUnauthorized, "invalid role")
				return
			}

			// Extract user ID
			var userID pgtype.UUID
			if err := userID.Scan(claims["user_id"]); err != nil {
				WriteError(w, http.StatusUnauthorized, "invalid user_id")
				return
			}

			// Store in context
			ctx := context.WithValue(r.Context(), ContextUserID, userID)
			ctx = context.WithValue(ctx, ContextRole, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetRole(r.Context())
		if !ok || role != "admin" {
			WriteError(w, http.StatusForbidden, "admin only")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SupplierOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetRole(r.Context())
		if !ok || (role != "supplier" && role != "admin") {
			WriteError(w, http.StatusForbidden, "supplier only")
			return
		}
		next.ServeHTTP(w, r)
	})
}
