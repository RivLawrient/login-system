package errs

import "errors"

var (
	ErrEmailUsed            = errors.New("email already used")
	ErrFailedCreateData     = errors.New("failed create database")
	ErrInternal             = errors.New("internal server error")
	ErrInvalidType          = errors.New("invalid type")
	ErrDataNotFound         = errors.New("data is not found")
	ErrInvalidEmailPassword = errors.New("email or password is invalid")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidAccessToken   = errors.New("invalid access token")
)
