package repo

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrBrokenData = errors.New("broken data in storage")
)
