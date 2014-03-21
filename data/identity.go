package data

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Identity struct {
	Id int64
}

func CreateIdentity(tx *sqlx.Tx, identity *Identity) error {
	const (
		Q = `INSERT INTO identities DEFAULT VALUES RETURNING *;`
	)

	return tx.Get(identity, Q)
}

func GetIdentity(q sqlx.Queryer, id int64) (*Identity, error) {
	const (
		Q = `SELECT * FROM identities WHERE id = $1 LIMIT 1;`
	)

	identity := &Identity{}
	err := sqlx.Get(q, identity, Q, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return identity, nil
}
