package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sheenobu/cm.go"
	"github.com/sheenobu/cm.go/sql"
)

// _Frameworks defines the collection for the Framework entity
type _Frameworks struct {
	cm.Collection
	Id          cm.ValueColumn
	Name        cm.ValueColumn
	Description cm.ValueColumn
	Url         cm.ValueColumn
}

// The instance of Frameworks used by the API consumers
var Frameworks *_Frameworks

// initialize the system database
func initDb() {

	// we are sqlite3 based so far
	db, _ := sqlx.Connect("sqlite3", "frameworks.db")

	// We initialize the collection using a SQL based backend
	Frameworks = &_Frameworks{
		Collection:  sql.New(db, "frameworks"),
		Id:          sql.Integer("id", 10).PrimaryKey(),
		Name:        sql.Varchar("name", 100).NotNull(),
		Description: sql.Varchar("description", 100).NotNull(),
		Url:         sql.Varchar("url", 100).NotNull(),
	}

	// Init performs various operations to analyze columns, create initial schema
	_ = Frameworks.Init(Frameworks) //TODO: check error
}
