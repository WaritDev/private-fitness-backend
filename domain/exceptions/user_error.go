package exceptions

import "errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)