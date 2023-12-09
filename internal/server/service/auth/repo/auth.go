package repo

import (
	"arch/internal/server/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

const loginQuery = `
select u.role_id permission
from "user" u
where u.login = $1
  and u.password = $2
`

func (r *Repo) Login(ctx context.Context, credentials entity.UserCredentials) (int, error) {
	var u userInfo
	if err := r.db.GetContext(ctx, &u, loginQuery, credentials.Login, credentials.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotAllowed
		}

		return 0, fmt.Errorf("r.db.Get: %w", err)
	}

	return u.Permission, nil
}

const checkPermissionQuery = `
call public.check_user_role($1, $2, $3)
`

func (r *Repo) CheckPermission(ctx context.Context, permission int, resource entity.ResourceInfo) error {
	if _, err := r.db.QueryxContext(ctx, checkPermissionQuery, resource.Name, resource.Method, permission); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Message {
			case sqlErrorNotAllowedRole:
				return ErrNotAllowed
			case sqlErrorResourceNotFound:
				return nil
			}
		}
		return fmt.Errorf("r.db.Exec: %w", err)
	}

	return nil
}
