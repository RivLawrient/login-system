package auth

import (
	"context"
	"database/sql"

	"github.com/RivLawrient/login-system/backend/internal/apps/domain/entity"
	"github.com/RivLawrient/login-system/backend/internal/apps/domain/repository"
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

func (s *Service) Register() {
	ctx := context.Background()
	tx, err := s.DB.Begin()
	if err != nil {
		panic(err)
	}

	s.UserRepo.Create(tx, ctx, &entity.User{})

	tx.Commit()
}
