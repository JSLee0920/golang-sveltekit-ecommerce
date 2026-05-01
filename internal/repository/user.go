package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	rdb     *redis.Client
}

func NewUserRepository(db *pgxpool.Pool, queries *generated.Queries, rdb *redis.Client) *UserRepository {
	return &UserRepository{db: db, queries: queries, rdb: rdb}
}

func userIDCacheKey(id pgtype.UUID) string {
	return fmt.Sprintf("user:id:%x", id.Bytes)
}

func (r *UserRepository) GetByID(ctx context.Context, id pgtype.UUID) (*generated.User, error) {
	cacheKey := userIDCacheKey(id)

	cached, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var user generated.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(user)
	r.rdb.Set(ctx, cacheKey, data, 30*time.Minute)

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*generated.User, error) {
	cacheKey := "user:email:" + email

	cached, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var user generated.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	// cache miss — query DB
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(user)
	r.rdb.Set(ctx, cacheKey, data, 30*time.Minute)

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

	// invalidate cache
	r.rdb.Del(ctx, userIDCacheKey(user.ID))
	r.rdb.Del(ctx, "user:email:"+user.Email)

	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	// get user first to invalidate email cache
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.queries.DeleteUser(ctx, id); err != nil {
		return err
	}

	// invalidate cache
	r.rdb.Del(ctx, userIDCacheKey(id))
	r.rdb.Del(ctx, "user:email:"+user.Email)

	return nil
}
