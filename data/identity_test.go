package data

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIdentity(t *testing.T) {
	Convey("with a database", t, func() {
		Convey("create identity", func() {
			db := sqlx.MustConnect("postgres", "postgres://localhost:15432/auth_sh?sslmode=disable")
			tx := db.MustBegin()

			defer db.Close()
			defer tx.Rollback()

			identity := &Identity{}
			err := CreateIdentity(tx, identity)
			So(err, ShouldBeNil)
			So(identity.Id, ShouldBeGreaterThan, 0)
		})
	})
}
