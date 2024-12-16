package storage

import "errors"

var (
	UserExists   = errors.New("user already exists")
	UserNotFound = errors.New("user not found")
)
