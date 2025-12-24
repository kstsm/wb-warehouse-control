package apperrors

import "errors"

var (
	ErrItemNotFound      = errors.New("item not found")
	ErrEmptyDate         = errors.New("empty date string")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrRoleMismatch      = errors.New("user role mismatch")
)
