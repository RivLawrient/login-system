package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/RivLawrient/login-system/backend/internal/apps/domain/entity"
	"github.com/RivLawrient/login-system/backend/internal/apps/domain/repository"
	"github.com/RivLawrient/login-system/backend/internal/errs"
	"github.com/RivLawrient/login-system/backend/internal/helper"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB       *sql.DB
	Redis    *redis.Client
	UserRepo *repository.UserRepo
}

func NewService(db *sql.DB, redis *redis.Client, userRepo *repository.UserRepo) *Service {
	return &Service{
		DB:       db,
		Redis:    redis,
		UserRepo: userRepo,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (*entity.User, *string, *string, error) {
	tx, err := s.DB.Begin()
	defer tx.Rollback()

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, nil, err
	}

	user := entity.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(pass),
	}

	err = s.UserRepo.Create(tx, ctx, &user)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, err
	}

	jwt, err := helper.GenerateJWT(user.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	refreshToken := uuid.NewString()
	key := "refresh_token:" + refreshToken
	if err := s.Redis.Set(
		ctx,
		key,
		user.ID,
		7*24*time.Hour,
	).Err(); err != nil {
		return nil, nil, nil, err
	}

	return &user, &jwt, &refreshToken, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*entity.User, *string, *string, error) {
	tx, _ := s.DB.Begin()
	defer tx.Rollback()

	user := entity.User{}
	if err := s.UserRepo.GetByEmail(tx, ctx, email, &user); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			return nil, nil, nil, errs.ErrInvalidEmailPassword
		}
		return nil, nil, nil, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, nil, errs.ErrInvalidEmailPassword
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, err
	}

	jwt, err := helper.GenerateJWT(user.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	refreshToken := uuid.NewString()
	key := "refresh_token:" + refreshToken
	if err := s.Redis.Set(
		ctx,
		key,
		user.ID,
		7*24*time.Hour,
	).Err(); err != nil {
		return nil, nil, nil, err
	}

	return &user, &jwt, &refreshToken, nil
}

func (s *Service) RefreshAccessToken(ctx context.Context, refreshToken string) (*string, *string, error) {
	key := "refresh_token:" + refreshToken
	userID, err := s.Redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil, errs.ErrInvalidRefreshToken
		}
		return nil, nil, err
	}

	// generate new access token
	accessToken, err := helper.GenerateJWT(userID)
	if err != nil {
		return nil, nil, err
	}

	// rotate refresh token
	newRefresh := uuid.NewString()
	newKey := "refresh_token:" + newRefresh
	if err := s.Redis.Set(ctx, newKey, userID, 7*24*time.Hour).Err(); err != nil {
		return nil, nil, err
	}

	// delete old token (best-effort)
	_ = s.Redis.Del(ctx, key).Err()

	return &accessToken, &newRefresh, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*entity.User, error) {
	tx, _ := s.DB.Begin()
	defer tx.Rollback()

	var user entity.User
	if err := s.UserRepo.GetByID(tx, ctx, userID, &user); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}
