package errs

import "errors"

var (
	ErrEmailUsed        = errors.New("email already used")
	ErrFailedCreateData = errors.New("failed create database")
	ErrInternal         = errors.New("internal server error")
)
