package repo

import (
	"arch/internal/server/entity"
	"errors"
)

const (
	sqlErrorNotAllowedRole   = "not allowed role"
	sqlErrorResourceNotFound = "resource not found"
)

var (
	ErrNotAllowed = errors.New("not allowed")
)

type userCredentials struct {
	ID         int    `db:"id"`
	Login      string `db:"login"`
	Password   string `db:"password"`
	Permission int    `db:"permission"`
}

func (u *userCredentials) serviceModel() entity.UserCredentials {
	return entity.UserCredentials{
		Login:      u.Login,
		Password:   u.Password,
		Permission: u.Permission,
	}
}
