package repository

import (
	"context"
	"testing"

	"github.com/RivLawrient/login-system/backend/internal/apps/domain/entity"
	"github.com/RivLawrient/login-system/backend/internal/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	userData := &entity.User{
		ID:       uuid.NewString(),
		Email:    "duplicate@test.com",
		Password: "hashedpassword",
	}

	tx1, _ := db.Begin()
	defer tx1.Rollback()

	err := repo.Create(tx1, ctx, userData)
	assert.NoError(t, err)
}

func TestUserRepo_Create_DuplicateEmail(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	user1 := &entity.User{
		ID:       uuid.NewString(),
		Email:    "duplicate2@test.com",
		Password: "hashedpassword",
	}

	user2 := &entity.User{
		ID:       uuid.NewString(),
		Email:    "duplicate2@test.com",
		Password: "anotherpassword",
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := repo.Create(tx, ctx, user1)
	assert.NoError(t, err)

	err = repo.Create(tx, ctx, user2)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrEmailUsed.Error())
}

func TestUserRepo_GetByID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	userData := &entity.User{
		ID:       uuid.NewString(),
		Email:    "getbyid@test.com",
		Password: "hashedpassword",
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := repo.Create(tx, ctx, userData)
	assert.NoError(t, err)

	var got entity.User
	err = repo.GetByID(tx, ctx, userData.ID, &got)
	assert.NoError(t, err)
	assert.Equal(t, userData.ID, got.ID)
	assert.Equal(t, userData.Email, got.Email)
	assert.Equal(t, userData.Password, got.Password)
}

func TestUserRepo_GetByID_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	tx, _ := db.Begin()
	defer tx.Rollback()

	var got entity.User
	err := repo.GetByID(tx, ctx, uuid.NewString(), &got)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrDataNotFound)
}

func TestUserRepo_GetByID_InvalidID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	tx, _ := db.Begin()
	defer tx.Rollback()

	var got entity.User
	err := repo.GetByID(tx, ctx, "not-a-uuid", &got)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrInvalidType)
}

func TestUserRepo_GetByEmail(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	userData := &entity.User{
		ID:       uuid.NewString(),
		Email:    "getbyemail@test.com",
		Password: "hashedpassword",
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := repo.Create(tx, ctx, userData)
	assert.NoError(t, err)

	var got entity.User
	err = repo.GetByEmail(tx, ctx, userData.Email, &got)
	assert.NoError(t, err)
	assert.Equal(t, userData.ID, got.ID)
	assert.Equal(t, userData.Email, got.Email)
	assert.Equal(t, userData.Password, got.Password)
}

func TestUserRepo_GetByEmail_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	tx, _ := db.Begin()
	defer tx.Rollback()

	var got entity.User
	err := repo.GetByEmail(tx, ctx, "notfound@example.com", &got)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrDataNotFound)
}
