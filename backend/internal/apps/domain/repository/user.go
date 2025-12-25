package repository

import (
	"context"
	"database/sql"

	"github.com/RivLawrient/login-system/backend/internal/apps/domain/entity"
	"github.com/RivLawrient/login-system/backend/internal/errs"
	"github.com/lib/pq"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(db *sql.Tx, ctx context.Context, data *entity.User) error {
	query := `
		INSERT INTO users(id, email, password) 
		VALUES($1, $2, $3)
	`

	result, err := db.ExecContext(ctx, query, data.ID, data.Email, data.Password)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return errs.ErrEmailUsed
			}
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrFailedCreateData
	}

	return err
}
