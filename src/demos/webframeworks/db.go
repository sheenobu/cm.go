package main

import (
	"cm"
	"cm/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
		Id:          sql.Column("id", "integer primary key AUTOINCREMENT"), // Should be abstracted out to a sql.PrimaryKeyInteger
		Name:        sql.Column("name", "varchar(100) not null"),           // Should be abstracted out to a sql.NotNullable(sql.VarChar(100))
		Description: sql.Column("description", "varchar(100) not null"),    // Should be abstracted out to a sql.NotNullable(sql.VarChar(100))
		Url:         sql.Column("url", "varchar(100) not null"),            // Should be abstracted out to a sql.NotNullable(sql.VarChar(100))
	}

	// Init performs various operations to analyze columns, create initial schema
	_ = Frameworks.Init(Frameworks) //TODO: check error
}
