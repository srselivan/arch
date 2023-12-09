package repo

import (
	"errors"
)

const (
	sqlErrorNotAllowedRole   = "not allowed role"
	sqlErrorResourceNotFound = "resource not found"
)

var (
	ErrNotAllowed = errors.New("not allowed")
)

type userInfo struct {
	Permission int `db:"permission"`
}
