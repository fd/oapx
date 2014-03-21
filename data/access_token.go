package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type AccessToken struct {
	Id            int64
	IdentityId    int64          `db:"identity_id"`
	ApplicationId int64          `db:"application_id"`
	AccessToken   string         `db:"access_token"`
	RefreshToken  sql.NullString `db:"refresh_token"`
	ExpiresIn     int32          `db:"expires_in"`
	CreatedAt     time.Time      `db:"created_at"`
	Scope         string
	RedirectURI   string `db:"redirect_uri"`
}

func CreateAccessToken(q sqlx.Queryer, access_token *AccessToken) error {
	const (
		Q = `INSERT INTO access_tokens (identity_id, application_id, access_token, refresh_token, expires_in, scope, redirect_uri) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;`
	)

	if access_token.RefreshToken.String != "" {
		access_token.RefreshToken.Valid = true
	}

	return sqlx.Get(q, access_token, Q,
		access_token.IdentityId,
		access_token.ApplicationId,
		access_token.AccessToken,
		access_token.RefreshToken,
		access_token.ExpiresIn,
		access_token.Scope,
		access_token.RedirectURI,
	)
}

func GetAccessTokenWithAccessToken(q sqlx.Queryer, token string) (*AccessToken, error) {
	const (
		Q = `SELECT * FROM access_tokens WHERE access_token = $1 LIMIT 1;`
	)

	access_token := &AccessToken{}
	err := sqlx.Get(q, access_token, Q, token)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return access_token, nil
}

func GetAccessTokenWithRefreshToken(q sqlx.Queryer, token string) (*AccessToken, error) {
	const (
		Q = `SELECT * FROM access_tokens WHERE refresh_token = $1 LIMIT 1;`
	)

	access_token := &AccessToken{}
	err := sqlx.Get(q, access_token, Q, token)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return access_token, nil
}

func DestroyAccessTokenWithAccessToken(e sqlx.Execer, token string) error {
	const (
		Q = `DELETE FROM access_tokens WHERE access_token = $1;`
	)

	_, err := e.Exec(Q, token)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func DestroyAccessTokenWithRefreshToken(e sqlx.Execer, token string) error {
	const (
		Q = `DELETE FROM access_tokens WHERE refresh_token = $1;`
	)

	_, err := e.Exec(Q, token)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
