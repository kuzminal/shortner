package linkit

import "errors"

var (
	ErrExists   = errors.New("already exists")
	ErrNotExist = errors.New("does not exist")
	ErrInternal = errors.New("internal error: please try again later or contact support")
)
