package project

import "errors"

var (
	ErrNotFound      = errors.New("project not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrForbidden     = errors.New("forbidden")
	ErrAlreadyMember = errors.New("already member")
)
