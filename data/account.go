package data

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/fd/oauth2-proxy/providers"
)

type Account struct {
	Id         int64
	IdentityId int64 `db:"identity_id"`

	ProviderId int64  `db:"provider_id"` // reference to provider info
	RemoteId   string `db:"remote_id"`   // id of the account with the provider

	Name    string
	Email   string
	Picture string

	RawProfile []byte `db:"raw_profile"`
	RawToken   []byte `db:"raw_token"`

	provider *providers.Provider
}

func CreateAccount(tx *sqlx.Tx, account *Account) error {
	const (
		Q = `INSERT INTO accounts (identity_id, remote_id, name, email, picture, raw_profile, raw_token) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;`
	)

	return tx.Get(account, Q,
		account.IdentityId,
		account.RemoteId,
		account.Name,
		account.Email,
		account.Picture,
		account.RawProfile,
		account.RawToken,
	)
}

func GetAccountWithRemoteId(tx *sqlx.Tx, remote_id string) (*Account, error) {
	const (
		Q = `SELECT * FROM accounts WHERE remote_id = $1 LIMIT 1;`
	)

	account := &Account{}
	err := tx.Get(account, Q, remote_id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

func GetAccountsWithIdentityId(tx *sqlx.Tx, identity_id int64) ([]*Account, error) {
	const (
		Q = `SELECT * FROM accounts WHERE identity_id = $1 ORDER BY remote_id;`
	)

	accounts := []*Account{}
	err := tx.Select(&accounts, Q, identity_id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func UpdateAccount(tx *sqlx.Tx, account *Account) error {
	const (
		Q = `UPDATE accounts SET name = $2, email = $3, picture = $4, raw_profile = $5, raw_token = $6 WHERE id = $1;`
	)

	_, err := tx.Exec(Q,
		account.Id,
		account.Name,
		account.Email,
		account.Picture,
		account.RawProfile,
		account.RawToken,
	)
	return err
}

func (a *Account) ProviderCode() string {
	idx := strings.IndexRune(a.RemoteId, ':')
	return a.RemoteId[:idx]
}

func (a *Account) Provider() *providers.Provider {
	if a.provider != nil {
		return a.provider
	}

	p, _ := providers.New(a.ProviderCode(), "", "")
	a.provider = p
	return p

}
