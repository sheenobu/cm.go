<!-- DO NOT EDIT Generated via template -->
<!-- vim: syntax=markdown
-->
# cm.go

Status: not sure... making it open source for now.

persistent collection management for golang, via declarative structs.

A few of the benifits of this approach are:

Separate model struct from persistence logic via a 'Collection structure':

```go
import cm "github.com/sheenobu/cm.go"

var Albums AlbumsCollection

type AlbumsCollection struct {
	cm.Collection
	ID       cm.ValueColumn
	Artist   cm.ValueColumn
	Name     cm.ValueColumn
	Year     cm.ValueColumn
	Explicit cm.ValueColumn
}
```

DSL for operating on these structures:

```go
// Querying for all released in the year 2011
var albumList []albums.Album
err = Albums.Filter(Albums.Year.Eq(2011)).List(ctx, &albumList)
```

## Get

	go get github.com/sheenobu/cm.go
	go get github.com/sheenobu/cm.go/sql

## Usage

albums/albums.go:

```go
import cm "github.com/sheenobu/cm.go"

// Album is the model for music album database
type Album struct {
	ID       *int
	Artist   string
	Name     string
	Explicit bool
	Year     int64
}

// AlbumsCollection defines the columns and operations for
// the Album model
type AlbumsCollection struct {
	cm.Collection
	ID       cm.ValueColumn
	Artist   cm.ValueColumn
	Name     cm.ValueColumn
	Year     cm.ValueColumn
	Explicit cm.ValueColumn
}
```

albums/albums\_sqlite.go:

```go
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
```

main.go:

```go
import (
	"context"
	"github.com/sheenobu/cm.go/doc/albums"
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
 * Automatic column generation on insert
 * (Optional) Transactions

## TODO

 * [DONE] Nullable column types
 * [DONE] Remove/Replace build tool to support importing from other projects.
 * [DONE] sql.Varchar
 * [DONE] In-place updates: c.Edit(c.MyColumn.Set("value"))
 * [DONE] Point to $GOPATH/../www for demo file root
 * [DONE] sql.Integer
 * [DONE] PrimaryKey, AutoIncrement builder functions
 * [DONE] Automatic ID generation via function
 * [DONE] Transactions
 * sql.DateTime
 * Raw Query API
 * Remove reflection code.
 * Updating by entity (Not sure if we want to support this)
 * Deleting by entity
 * Database Versioning
 * Relations / SQL joins
 * Caching
 * In-place updates:
    * c.Edit(c.MyColumn.Append("\_appended\_string") - TODO
	* c.Edit(c.MyIntegerColumn.Add(1)) - TODO
	* c.Edit(c.MyColumn.SetFunc(fn)) - TODO

