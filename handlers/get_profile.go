package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"

	"github.com/fd/oauth2-proxy/data"
)

func GET_profile(rw http.ResponseWriter, r *http.Request, identity *data.Identity,
	render render.Render, db *sqlx.DB) {
	var (
		tx      = db.MustBegin()
		success bool
		page    struct {
			Accounts          []*data.Account
			OwnedApplications []*data.Application
		}
	)

	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	accounts, err := data.GetAccountsWithIdentityId(tx, identity.Id)
	if err != nil {
		panic(err)
	}

	owned_applications, err := data.GetApplicationsOwnedByIdentity(tx, identity.Id)
	if err != nil {
		panic(err)
	}

	page.Accounts = accounts
	page.OwnedApplications = owned_applications
	render.HTML(200, "profile", &page)
}
