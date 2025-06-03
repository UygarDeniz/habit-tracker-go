package apperrors

import "errors"

var ErrNotFound = errors.New("requested resource not found")

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrInvalidInput = errors.New("invalid input provided")

var ErrDatabase = errors.New("database operation failed")
