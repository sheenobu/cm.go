package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"github.com/sheenobu/cm.go"
	"github.com/sheenobu/cm.go/sql"
)

// _Frameworks defines the collection for the Framework entity
type _Frameworks struct {
	cm.Collection
	ID          cm.ValueColumn
	Name        cm.ValueColumn
	Description cm.ValueColumn
	URL         cm.ValueColumn
}

func uuidGen() interface{} {
	return uuid.NewV1().String()
}

// Frameworks is the database attached instance
// of the frameworks collection
var Frameworks *_Frameworks

// initialize the system database
func initDB() {

	// we are sqlite3 based so far
	db, err := sqlx.Connect("sqlite3", "frameworks.db")
	if err != nil {
		panic(err)
	}

	// We initialize the collection using a SQL based backend
	Frameworks = &_Frameworks{
		Collection:  sql.New(db, "frameworks"),
		ID:          sql.Varchar("id", 32).PrimaryKey().FromFunction(uuidGen),
		Name:        sql.Varchar("name", 100).NotNull(),
		Description: sql.Varchar("description", 100).NotNull(),
		URL:         sql.Varchar("url", 100).NotNull(),
	}

	// Init performs various operations to analyze columns,
	// create initial schema
	err = Frameworks.Init(Frameworks)
	if err != nil {
		i := sql.GetErrorCode(err)
		switch i {
		case sql.TableAlreadyExists:
		default:
			panic(err)
		}
	}
}
