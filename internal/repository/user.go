package repository

import (
	"context"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db      *pgxpool.Pool
	queries *generated.Queries
}

func NewUserRepository(db *pgxpool.Pool, queries *generated.Queries) *UserRepository {
	return &UserRepository{db: db, queries: queries}
}

func (r *UserRepository) GetByID(ctx context.Context, id pgtype.UUID) (*generated.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*generated.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, params generated.CreateUserParams) (*generated.User, error) {
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, params generated.UpdateUserParams) (*generated.User, error) {
	user, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}
