package data

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccount(t *testing.T) {
	Convey("with a database", t, func() {
		Convey("create account", func() {
			db := sqlx.MustConnect("postgres", "postgres://localhost:15432/auth_sh?sslmode=disable")
			tx := db.MustBegin()

			defer db.Close()
			defer tx.Rollback()

			identity := &Identity{}
			err := CreateIdentity(tx, identity)

			account := &Account{
				IdentityId: identity.Id,
				RemoteId:   "fb:145478142",
				Name:       "Simon Menke",
				Email:      "simon.menke@gmail.com",
				Picture:    "",
				RawProfile: []byte("{}"),
				RawToken:   []byte("{}"),
			}
			err = CreateAccount(tx, account)
			So(err, ShouldBeNil)
			So(account.Id, ShouldBeGreaterThan, 0)
		})
	})
}
