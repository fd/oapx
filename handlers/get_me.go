package handlers

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"

	"github.com/fd/oauth2-proxy/data"
)

func GET_me(identity *data.Identity, render render.Render, db *sqlx.DB) {
	type account struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email,omitempty"`
		Picture string `json:"picture,omitempty"`
	}

	var (
		tx      = db.MustBegin()
		success bool
		err     error
		resp    struct {
			Id       string     `json:"id"`
			Accounts []*account `json:"accounts"`
		}
	)

	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	resp.Id = strconv.FormatInt(identity.Id, 10)

	accounts, err := data.GetAccountsWithIdentityId(tx, identity.Id)
	if err != nil {
		panic(err)
	}

	for _, acc := range accounts {
		resp.Accounts = append(resp.Accounts, &account{
			Id:      acc.RemoteId,
			Name:    acc.Name,
			Email:   acc.Email,
			Picture: acc.Picture,
		})
	}

	success = true
	render.JSON(200, &resp)
}
