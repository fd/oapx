package handlers

import (
	"net/http"

	"github.com/martini-contrib/sessions"
)

func GET_logout(sess sessions.Session, rw http.ResponseWriter, r *http.Request) {
	sess.Delete("identity_id")
	http.Redirect(rw, r, "/", http.StatusFound)
}
