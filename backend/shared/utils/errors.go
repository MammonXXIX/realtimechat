package utils

import "errors"

var (
	ErrNotFound  = errors.New("Not Found")
	ErrDuplicate = errors.New("Duplicate")
)
