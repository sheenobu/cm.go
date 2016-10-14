// +build ignore

package albums

import (
	"github.com/jmoiron/sqlx"
	"github.com/sheenobu/cm.go/sql"
)

// Collection is the albums collection
var Collection *_Albums

func init() {
	db, _ := sqlx.Connect("sqlite3", "albums.db")
	Collection = &_Albums{
		Collection: New(db, "ALBUMS"),
		//ID:		sql.Varchar("id", 32).PrimaryKey().FromFunction(uuidGen)

		ID:       sql.Integer().PrimaryKey(),
		Artist:   sql.Varchar("artist", 100).NotNull(),
		Name:     sql.Varchar("name", 100).NotNull(),
		Year:     sql.Column("year", "number not null"),
		Explicit: sql.Column("explicit", "bool not null default false"),
	}
	err := Collection.Init(Collection)
}
