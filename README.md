# cm.go

persistent collection management for golang

## Usage

albums/albums.go:

	package albums

	import (
		"cm"
		"cm/sql"
	)


	// Album defines the model for the music album
	type Album struct {
		Id string // BUG CM01
		Artist 	   string
		Name 	   string
		Explicit   bool
		Year       int64
	}

	// _Albums is the collection for the model.
	type _Albums struct {
		cm.Collection
		Id 		   cm.ValueColumn
		Artist 	   cm.ValueColumn
		Name 	   cm.ValueColumn
		Year       cm.ValueColumn
		Explicit   cm.ValueColumn
	}

	var Collection *_Albums

	func init() {
		db, _ := sqlx.Connect("sqlite3", "albums.db")
		Collection = &_Albums{
			Collection: New(db, "ALBUMS"),
			Id: sql.Column("id", "integer primary key"),
			Artist: 	    sql.Column("artist", "varchar(100) not null"),
			Name:      		sql.Column("name", "varchar(100) not null"),
			Year:       	sql.Column("year", "number not null"),
			Explicit:  		sql.Column("explicit", "bool not null default false"),
		}
		err := Collection.Init(Collection)
	}

main.go:

	package main

	import (
		"albums"
	)

	// optional shortcut for our example code
	var Albums *albums._Albums = albums.Collection

	func main() {
		// for potential transactions
		ctx := context.Background()

		// insertion
		err = Albums.Insert(ctx, &albums.Album{
			Artist: "Childish Gambino",
			Name:   "Camp",
			Year:   2011,
			Explicit: true,
		})

		err = Albums.Insert(ctx, &albums.Album{
			Artist: "Sleater-Kinney",
			Name:   "No Cities To Love",
			Year:   2015,
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
		err = Albums.Filter(Albums.Id.Eq(0)).Delete(ctx)
	}

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

 * Raw Query API
 * Remove reflection code.
 * Remove/Replace build tool to support importing from other projects.
 * Higher level column types (sql.VarChar, sql.Integer, sql.DateTime)
 * Nullable column types
 * Automatic ID generation
 * Updating by entity (Not sure if we want to support this)
 * Deleting by entity
 * Database Versioning
 * Relations / SQL joins
 * Caching
 * Transactions
 * In-place updates:
	* c.Edit(c.MyColumn.Set("value")) - DONE
	* c.Edit(c.MyColumn.Append("\_appended\_string") - TODO
	* c.Edit(c.MyIntegerColumn.Add(1)) - TODO
	* c.Edit(c.MyColumn.SetFunc(fn)) - TODO

## Bugs

 * CM01 - AUTO INCREMENT columns must be string type, otherwise the system attempts to insert a 0 for every row.
