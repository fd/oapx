package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Authorization struct {
	Id            int64
	IdentityId    int64 `db:"identity_id"`
	ApplicationId int64 `db:"application_id"`
	Code          string
	ExpiresIn     int32     `db:"expires_in"`
	CreatedAt     time.Time `db:"created_at"`
	State         string
	Scope         string
	RedirectURI   string `db:"redirect_uri"`
}

func CreateAuthorization(q sqlx.Queryer, authorization *Authorization) error {
	const (
		Q = `INSERT INTO authorizations (identity_id, application_id, code, expires_in, state, scope, redirect_uri) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;`
	)

	return sqlx.Get(q, authorization, Q,
		authorization.IdentityId,
		authorization.ApplicationId,
		authorization.Code,
		authorization.ExpiresIn,
		authorization.State,
		authorization.Scope,
		authorization.RedirectURI,
	)
}

func GetAuthorizationWithCode(q sqlx.Queryer, code string) (*Authorization, error) {
	const (
		Q = `SELECT * FROM authorizations WHERE code = $1 LIMIT 1;`
	)

	authorization := &Authorization{}
	err := sqlx.Get(q, authorization, Q, code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return authorization, nil
}

func DestroyAuthorizationWithCode(e sqlx.Execer, code string) error {
	const (
		Q = `DELETE FROM authorizations WHERE code = $1;`
	)

	_, err := e.Exec(Q, code)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
