package repository

import (
	"database/sql"
	"testing"

	"github.com/RivLawrient/login-system/backend/config"
)

func SetupTestDB(t *testing.T) *sql.DB {
	config.LoadEnv("../../../../.env")
	db := config.GetConnection()

	return db
}
