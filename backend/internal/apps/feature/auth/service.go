package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/RivLawrient/login-system/backend/internal/apps/domain/entity"
	"github.com/RivLawrient/login-system/backend/internal/apps/domain/repository"
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

	//code
	//get user
	//compare password
	//genereate jwt
	//genereate refresh

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, err
	}

	return nil, nil, nil, nil
}
