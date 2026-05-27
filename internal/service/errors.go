package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrTicketNotFound     = errors.New("ticket not found")
	ErrAccessDenied       = errors.New("access denied")
)
