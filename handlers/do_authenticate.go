package handlers

import (
	"net/http"
	"strings"

	"github.com/codegangsta/martini"
	"github.com/jmoiron/sqlx"
	"github.com/martini-contrib/sessions"

	"github.com/fd/oauth2-proxy/data"
)

func MustAuthenticate() martini.Handler {
	return must_authenticate
}

func must_authenticate(c martini.Context, sess sessions.Session, db *sqlx.DB, r *http.Request) {
	identity := ActiveIdentity(c)

	if identity != nil {
		return
	}

	if r.Header.Get("x-interactive") == "true" {
		sess.Delete("identity_id")
		c.Invoke(redirect_to("/login"))
	} else {
		c.Invoke(forbidden())
	}
}

func MayAuthenticate() martini.Handler {
	return may_authenticate
}

func may_authenticate(c martini.Context, sess sessions.Session, db *sqlx.DB, r *http.Request) {
	var (
		interactive = true
		token       string
		identity_id int64
		identity    *data.Identity
		err         error
	)

	// Attempt with Authorization header
	if v := r.Header.Get("Authorization"); v != "" {
		parts := strings.SplitN(v, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			interactive = false
			token = parts[1]
		}

		// Attempt with access_token parameter
	} else if v := r.URL.Query().Get("access_token"); v != "" {
		interactive = false
		token = v

		// Attempt with session.identity_id
	} else if id, ok := sess.Get("identity_id").(int64); ok {
		interactive = true
		identity_id = id
	}

	if token != "" {
		at, err := data.GetAccessTokenWithAccessToken(db, token)
		if err != nil {
			panic(err)
		}
		identity_id = at.IdentityId
	}

	if identity_id > 0 {
		identity, err = data.GetIdentity(db, identity_id)
		if err != nil {
			panic(err)
		}
	}

	if interactive {
		r.Header.Set("x-interactive", "true")
	}

	c.Map(identity)
}

func ActiveIdentity(c martini.Context) *data.Identity {
	var (
		identity *data.Identity
	)

	c.Invoke(MayAuthenticate())

	c.Invoke(func(i *data.Identity) { identity = i })

	return identity
}
