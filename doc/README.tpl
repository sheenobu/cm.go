<!-- vim: syntax=markdown
-->
# cm.go

persistent collection management for golang

## Get

	go get github.com/sheenobu/cm.go
	go get github.com/sheenobu/cm.go/sql

## Usage

albums/albums.go:

```
{{ shell "tail -n+5 ./albums/albums.go" }}
```

albums/albums\_sqlite.go:

```
{{ shell "tail -n+5 ./albums/albums_sqlite.go" }}
```

main.go:

```
{{ shell "tail -n+5 ./main.go" }}
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

## TODO

 * [DONE] Nullable column types
 * [DONE] Remove/Replace build tool to support importing from other projects.
 * [DONE] sql.Varchar
 * [DONE] In-place updates: c.Edit(c.MyColumn.Set("value"))
 * [DONE] Point to $GOPATH/../www for demo file root
 * [DONE] sql.Integer
 * [DONE] PrimaryKey, AutoIncrement builder functions
 * [DONE] Automatic ID generation via function
 * sql.DateTime
 * Raw Query API
 * Remove reflection code.
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

