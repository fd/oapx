package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/jmoiron/sqlx"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	"github.com/fd/oauth2-proxy/handlers"
	"github.com/fd/oauth2-proxy/oauth"
)

func main() {
	db := sqlx.MustConnect("postgres", "postgres://localhost:15432/auth_sh?sslmode=disable")
	defer db.Close()

	s := oauth.New(db)

	m := martini.Classic()
	m.Use(handlers.NoCache())
	m.Use(sessions.Sessions("auth-sh", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(render.Renderer(render.Options{Layout: "layout"}))
	m.Use(dump_request)
	m.Map(db)
	m.Map(s)

	m.Get("/login", handlers.GET_login)
	m.Get("/link", handlers.GET_link)
	m.Any("/authorize", handlers.GET_authorize)
	m.Any("/token", handlers.GET_token)
	m.Any("/info", handlers.GET_info)
	m.Get("/continue/:provider", handlers.GET_continue)
	m.Get("/logout", handlers.GET_logout)
	m.Get("/callback", handlers.GET_callback)

	m.Get("/", handlers.MayAuthenticate(), handlers.GET_home)
	m.Get("/me", handlers.MustAuthenticate(), handlers.GET_me)
	m.Get("/profile", handlers.MustAuthenticate(), handlers.GET_profile)
	m.Post("/applications", handlers.MustAuthenticate(), handlers.POST_application)

	m.Run()
}

func dump_request(r *http.Request, log *log.Logger) {
	for k, v := range r.URL.Query() {
		log.Printf("Q: %s: %s", k, v)
	}
	for k, v := range r.Header {
		log.Printf("H: %s: %s", k, v)
	}
}
