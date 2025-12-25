package auth

import (
	"context"
	"database/sql"

	"github.com/RivLawrient/login-system/backend/internal/apps/domain/entity"
	"github.com/RivLawrient/login-system/backend/internal/apps/domain/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB       *sql.DB
	UserRepo *repository.UserRepo
}

func NewService(db *sql.DB, userRepo *repository.UserRepo) *Service {
	return &Service{
		DB:       db,
		UserRepo: userRepo,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (*entity.User, error) {
	tx, err := s.DB.Begin()
	defer tx.Rollback()

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(pass),
	}

	err = s.UserRepo.Create(tx, ctx, &user)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}
