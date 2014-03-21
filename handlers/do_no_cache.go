package handlers

import (
	"net/http"

	"github.com/codegangsta/martini"
)

func NoCache() martini.Handler {
	return func(res http.ResponseWriter) {
		rw := res.(martini.ResponseWriter)
		rw.Before(func(martini.ResponseWriter) {
			rw.Header().Set("Cache-Control", "no-cache")
		})
	}
}
