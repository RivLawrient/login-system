package internal

import (
	"database/sql"

	"github.com/RivLawrient/login-system/backend/internal/apps/domain/repository"
	"github.com/RivLawrient/login-system/backend/internal/apps/feature/auth"
	"github.com/RivLawrient/login-system/backend/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type AppsConfig struct {
	DB       *sql.DB
	Redis    *redis.Client
	App      *gin.Engine
	Validate *validator.Validate
}

func Apps(a *AppsConfig) {
	userRepo := repository.NewUserRepo()

	authServ := auth.NewService(a.DB, a.Redis, userRepo)
	authHand := auth.NewHandler(authServ, a.Validate)

	route.Routes{
		App:         a.App,
		AuthHandler: authHand,
	}.SetupRouts()
}
