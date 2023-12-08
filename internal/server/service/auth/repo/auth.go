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

const credentialsQuery = `
select u.id,
       u.login,
       u.password,
       u.role_id permission
from "user" u
where u.login = $1
  and u.password = $2
`

func (r *Repo) Credentials(ctx context.Context, credentials entity.UserCredentials) (entity.UserCredentials, error) {
	var userCreds userCredentials
	if err := r.db.GetContext(ctx, &userCreds, credentialsQuery, credentials.Login, credentials.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.UserCredentials{}, ErrNotAllowed
		}

		return entity.UserCredentials{}, fmt.Errorf("r.db.Get: %w", err)
	}

	return userCreds.serviceModel(), nil
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
