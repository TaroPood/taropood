package domain

import "errors"

var (
	ErrNotFound      = errors.New("rule not found")
	ErrInvalidRule   = errors.New("invalid rule")
	ErrDuplicateName = errors.New("rule name already exists")
)
