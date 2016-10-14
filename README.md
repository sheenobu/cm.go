<!-- DO NOT EDIT Generated via template -->
# cm.go

persistent collection management for golang

## Get

	go get github.com/sheenobu/cm.go
	go get github.com/sheenobu/cm.go/sql

## Usage

albums/albums.go:

```
import (
	"github.com/sheenobu/cm.go"
)

// Album defines the model for the music album
type Album struct {
	ID       *int
	Artist   string
	Name     string
	Explicit bool
	Year     int64
}

// _Albums is the collection for the model.
type _Albums struct {
	cm.Collection
	ID       cm.ValueColumn
	Artist   cm.ValueColumn
	Name     cm.ValueColumn
	Year     cm.ValueColumn
	Explicit cm.ValueColumn
}
```

albums/albums\_sqlite.go:

```
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
		ID:         sql.Integer().PrimaryKey(),
		Artist:     sql.Varchar("artist", 100).NotNull(),
		Name:       sql.Varchar("name", 100).NotNull(),
		Year:       sql.Column("year", "number not null"),
		Explicit:   sql.Column("explicit", "bool not null default false"),
	}
	err := Collection.Init(Collection)
}
```

main.go:

```
import (
	"albums"
	"context"
)

// Albums is an optional shortcut for our example code
var Albums = albums.Collection

func main() {
	// for potential transactions
	ctx := context.Background()

	// insertion
	err = Albums.Insert(ctx, &albums.Album{
		Artist:   "Childish Gambino",
		Name:     "Camp",
		Year:     2011,
		Explicit: true,
	})

	err = Albums.Insert(ctx, &albums.Album{
		Artist:   "Sleater-Kinney",
		Name:     "No Cities To Love",
		Year:     2015,
		Explicit: false,
	})

	// Querying for all released in the year 2011
	var albumList []albums.Album
	err = Albums.Filter(Albums.Year.Eq(2011)).List(ctx, &albumList)

	// Updating all explicit albums to non-explicit
	filter := Albums.Explicit.Eq(true)
	setter := Albums.Explicit.Set(false)
	err = Albums.Filter(filter).Edit(setter).Update(ctx)

	// Deletion by ID
	err = Albums.Filter(Albums.ID.Eq(0)).Delete(ctx)
}
```

## Backends

 * sqlx 		- In Progress
 * database/sql - TODO
 * others		- TBD

## Implemented

 * Filtering
 * Fetching
 * Deleting via filtering
 * Updating via filtering
 * Simple schema generation
 * Inserting
 * Pagination

## TODO

 * [DONE] Nullable column types
 * [DONE] Remove/Replace build tool to support importing from other projects.
 * [DONE] sql.Varchar
 * [DONE] In-place updates: c.Edit(c.MyColumn.Set("value"))
 * [DONE] Point to $GOPATH/../www for demo file root
 * [DONE] sql.Integer
 * [DONE] PrimaryKey, AutoIncrement builder functions
 * sql.DateTime
 * Raw Query API
 * Remove reflection code.
 * Automatic ID generation via function
 * Updating by entity (Not sure if we want to support this)
 * Deleting by entity
 * Database Versioning
 * Relations / SQL joins
 * Caching
 * Transactions
 * In-place updates:
    * c.Edit(c.MyColumn.Append("\_appended\_string") - TODO
	* c.Edit(c.MyIntegerColumn.Add(1)) - TODO
	* c.Edit(c.MyColumn.SetFunc(fn)) - TODO

