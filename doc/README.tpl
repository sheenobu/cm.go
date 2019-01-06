<!-- vim: syntax=markdown
-->
# cm.go

Status: not sure... making it open source for now.

persistent collection management for golang, via declarative structs.

A few of the benifits of this approach are:

Separate model struct from persistence logic via a 'Collection structure':

```go
{{ shell "tail -n+5 ./albums/albums_single.go" }}
```

DSL for operating on these structures:

```go
{{ shell "tail -n+2 ./albums/albums_dsl.go" }}
```

## Get

	go get github.com/sheenobu/cm.go
	go get github.com/sheenobu/cm.go/sql

## Usage

albums/albums.go:

```go
{{ shell "tail -n+5 ./albums/albums.go" }}
```

albums/albums\_sqlite.go:

```go
{{ shell "tail -n+5 ./albums/albums_sqlite.go" }}
```

main.go:

```go
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

