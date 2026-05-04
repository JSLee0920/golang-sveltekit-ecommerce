package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/dto"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/middleware"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/response"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	svc *service.UserService
	cfg *config.Config
}

func NewUserHandler(svc *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{svc: svc, cfg: cfg}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var body dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, response.Failure[string]("invalid body"))
		return
	}
	if err := body.Validate(); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, response.Failure[string](err.Error()))
		return
	}

	user, err := h.svc.Create(r.Context(), body.Name, body.Email, body.Password, "customer")
	if err != nil {
		middleware.WriteJSON(w, http.StatusConflict, response.Failure[string](err.Error()))
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, response.Success(dto.FromUser(user)))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var body dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, response.Failure[string]("invalid body"))
		return
	}
	if err := body.Validate(); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, response.Failure[string](err.Error()))
		return
	}

	user, err := h.svc.GetByEmail(r.Context(), body.Email)
	if err != nil {
		middleware.WriteJSON(w, http.StatusUnauthorized, response.Failure[string]("invalid credentials"))
		return
	}
	if err := h.svc.VerifyPassword(user.Password, body.Password); err != nil {
		middleware.WriteJSON(w, http.StatusUnauthorized, response.Failure[string]("invalid credentials"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		middleware.WriteJSON(w, http.StatusInternalServerError, response.Failure[string]("failed to generate token"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    signed,
		HttpOnly: true,
		Secure:   h.cfg.AppEnv == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	middleware.WriteJSON(w, http.StatusOK, response.Success(dto.FromUser(user)))
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
	middleware.WriteJSON(w, http.StatusOK, response.Success("Logged out"))
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextUserID).(pgtype.UUID)
	if !ok {
		middleware.WriteJSON(w, http.StatusUnauthorized, response.Failure[string]("unauthorized"))
		return
	}

	user, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		middleware.WriteJSON(w, http.StatusNotFound, response.Failure[string]("user not found"))
		return
	}

	middleware.WriteJSON(w, http.StatusOK, response.Success(dto.FromUser(user)))
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	userID, ok := r.Context().Value(middleware.ContextUserID).(pgtype.UUID)
	if !ok {
		middleware.WriteJSON(w, http.StatusUnauthorized, response.Failure[string]("unauthorized"))
		return
	}

	var body dto.UpdateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, response.Failure[string]("invalid body"))
		return
	}
	if err := body.Validate(); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, response.Failure[string](err.Error()))
		return
	}

	user, err := h.svc.Update(r.Context(), userID, body.Name)
	if err != nil {
		middleware.WriteJSON(w, http.StatusInternalServerError, response.Failure[string]("failed to update user"))
		return
	}

	middleware.WriteJSON(w, http.StatusOK, response.Success(dto.FromUser(user)))
}
