package oauth

import (
	"github.com/RangelReale/osin"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *osin.Server {
	return osin.NewServer(config, &Storage{db})
}
