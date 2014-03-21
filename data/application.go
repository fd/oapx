package data

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"io"

	"github.com/jmoiron/sqlx"
)

type Application struct {
	Id      int64
	OwnerId int64 `db:"owner_id"`
	Name    string

	ClientId     string `db:"client_id"`
	ClientSecret string `db:"client_secret"`
	RedirectURI  string `db:"redirect_uri"`
}

func CreateApplication(tx *sqlx.Tx, application *Application) error {
	const (
		Q = `INSERT INTO applications (owner_id, name, client_id, client_secret, redirect_uri) VALUES ($1, $2, $3, $4, $5) RETURNING *;`
	)

	rand_a, err := make_rand(16)
	if err != nil {
		return err
	}

	rand_b, err := make_rand(32)
	if err != nil {
		return err
	}

	return tx.Get(application, Q,
		application.OwnerId,
		application.Name,
		rand_a,
		rand_b,
		application.RedirectURI,
	)
}

func GetApplicationWithId(q sqlx.Queryer, id int64) (*Application, error) {
	const (
		Q = `SELECT * FROM applications WHERE id = $1 LIMIT 1;`
	)

	application := &Application{}
	err := sqlx.Get(q, application, Q, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return application, nil
}

func GetApplicationWithClientId(q sqlx.Queryer, client_id string) (*Application, error) {
	const (
		Q = `SELECT * FROM applications WHERE client_id = $1 LIMIT 1;`
	)

	application := &Application{}
	err := sqlx.Get(q, application, Q, client_id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return application, nil
}

func GetApplicationsOwnedByIdentity(tx *sqlx.Tx, identity_id int64) ([]*Application, error) {
	const (
		Q = `SELECT * FROM applications WHERE owner_id = $1 ORDER BY name;`
	)

	applications := []*Application{}
	err := tx.Select(&applications, Q, identity_id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func make_rand(l int) (string, error) {
	buf := make([]byte, l/2)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
