package internal

import (
	"database/sql"

	"github.com/RivLawrient/login-system/backend/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppsConfig struct {
	DB       *sql.DB
	App      *gin.Engine
	Validate *validator.Validate
}

func Apps(a *AppsConfig) {

	route.Routes{App: a.App}.SetupRouts()
}
