package service

import (
	"context"
	"errors"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/db/generated"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(ctx context.Context, id pgtype.UUID) (*generated.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*generated.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) Create(ctx context.Context, name, email, password, role string) (*generated.User, error) {
	// check if email already exists
	existing, _ := s.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user, err := s.repo.Create(ctx, generated.CreateUserParams{
		Name:     name,
		Email:    email,
		Password: string(hash),
		Role:     role,
	})
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, id pgtype.UUID, name string) (*generated.User, error) {
	user, err := s.repo.Update(ctx, generated.UpdateUserParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return nil, errors.New("failed to update user")
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id pgtype.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) VerifyPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
