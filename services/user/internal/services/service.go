package service

import "errors"

var (
	ErrEmailExists        = errors.New("email is already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
