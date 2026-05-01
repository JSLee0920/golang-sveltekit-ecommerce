package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/db/generated"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/middleware"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/service"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	svc *service.UserService
	cfg *config.Config
}

type UserResponse struct {
	ID        pgtype.UUID        `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Role      string             `json:"role"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func toUserResponse(user *generated.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserHandler(svc *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{svc: svc, cfg: cfg}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "name, email, and password are required"})
		return
	}

	user, err := h.svc.Create(r.Context(), body.Name, body.Email, body.Password, "customer")
	if err != nil {
		middleware.JSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	middleware.JSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	if body.Email == "" || body.Password == "" {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
		return
	}

	user, err := h.svc.GetByEmail(r.Context(), body.Email)
	if err != nil {
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	if err := h.svc.VerifyPassword(user.Password, body.Password); err != nil {
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
		return
	}

	// set JWT as HTTP only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    signed,
		HttpOnly: true,
		Secure:   h.cfg.AppEnv == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	middleware.JSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   h.cfg.AppEnv == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
	})
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(pgtype.UUID)

	user, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		middleware.JSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	middleware.JSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(pgtype.UUID)

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	if body.Name == "" {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	user, err := h.svc.Update(r.Context(), userID, body.Name)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	middleware.JSON(w, http.StatusOK, toUserResponse(user))
}
