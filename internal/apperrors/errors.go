package apperrors

import "errors"

var ErrNotFound = errors.New("requested resource not found")

var ErrInvalidInput = errors.New("invalid input provided")

var ErrDatabase = errors.New("database operation failed")

var ErrForbidden = errors.New("user is not authorized to perform this action")

var ErrAlreadyExists = errors.New("resource already exists")
