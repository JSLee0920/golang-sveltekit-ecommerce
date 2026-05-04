package dto

import (
	"errors"
	"time"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/db/generated"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	if r.Name == "" || r.Email == "" || r.Password == "" {
		return errors.New("name, email, and password are required")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if r.Email == "" || r.Password == "" {
		return errors.New("email and password are required")
	}
	return nil
}

type UpdateMeRequest struct {
	Name string `json:"name"`
}

func (r UpdateMeRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FromUser(u *generated.User) UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Time.Format(time.RFC3339),
	}
}
