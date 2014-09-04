package handlers

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/fd/oauth2-proxy/data"
)

func GET_home(c martini.Context, identity *data.Identity, render render.Render) {
	if identity == nil {
		c.Invoke(redirect_to("/login"))
	} else {
		c.Invoke(GET_profile)
	}
}
