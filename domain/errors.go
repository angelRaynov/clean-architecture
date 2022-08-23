package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("the requested item is not found")
	ErrConflict            = errors.New("this item already exist")
	ErrBadInput            = errors.New("invalid parameter")
)
