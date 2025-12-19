package task

import "errors"

var (
	ErrNotFound  = errors.New("task not found")
	ErrForbidden = errors.New("forbidden")
)
