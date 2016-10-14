package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sheenobu/cm.go/tx"
)

// shared interface for both sql transaction and
// sql databases. Could grow if you need another function
type dbI interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func getDb(ctx context.Context, db *sqlx.DB) dbI {
	t, ok := tx.Get(ctx)

	// only operate if we have an active transaction
	if !ok || t == nil || !t.Active() {

		// cast to sql transaction
		st, ok := t.(*Transaction)
		if !ok {
			return db
		}

		// get the internal transaction object
		return st.tx
	}

	return db
}
