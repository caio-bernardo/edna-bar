package types

import "errors"

var (
	ErrNotFound = errors.New("Not found")
	ErrInternalServer = errors.New("Internal error")
)
