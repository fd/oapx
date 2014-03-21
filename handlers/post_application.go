package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"

	"github.com/fd/oauth2-proxy/data"
)

func POST_application(rw http.ResponseWriter, r *http.Request, identity *data.Identity,
	render render.Render, db *sqlx.DB) {
	var (
		tx      = db.MustBegin()
		success bool
	)

	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	err = data.CreateApplication(tx, &data.Application{
		OwnerId:     identity.Id,
		Name:        r.PostForm.Get("application[name]"),
		RedirectURI: r.PostForm.Get("application[redirect_uri]"),
	})
	if err != nil {
		panic(err)
	}

	success = true
	http.Redirect(rw, r, "/", http.StatusFound)
}
