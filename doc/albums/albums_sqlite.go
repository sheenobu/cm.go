// +build ignore

package albums

import (
	"github.com/jmoiron/sqlx"
	"github.com/sheenobu/cm.go/sql"
)

// Collection is the database attached reference
// of the albums collection
var Collection *AlbumsCollection

func init() {
	db, _ := sqlx.Connect("sqlite3", "albums.db")

	// initialize the fields, declaring what they do and how they map
	Collection = &AlbumsCollection{
		Collection: sql.New(db, "ALBUMS"), // db connection, table name

		// if you want to use uuid's or some other generated key
		//ID:		sql.Varchar("id", 32).PrimaryKey().FromFunction(uuidGen)

		ID:       sql.Integer().PrimaryKey(),
		Artist:   sql.Varchar("artist", 100).NotNull(),
		Name:     sql.Varchar("name", 100).NotNull(),
		Year:     sql.Column("year", "number not null"),
		Explicit: sql.Column("explicit", "bool not null default false"),
	}

	// Init performs the heavy lifting of creating the tables,
	// pre-caching reflection results, building the querys.
	_ = Collection.Init(Collection)
}
