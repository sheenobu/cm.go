package sql

import (
	"github.com/jmoiron/sqlx"
	cm "github.com/sheenobu/cm.go"
)

// BeginTx begins a transaction
func (sql Table) BeginTx() (cm.Transaction, error) {
	t, err := sql.database.Beginx()

	if err != nil {
		return nil, err
	}

	return &Transaction{t}, nil
}

// Transaction is the sql-specific transaction struct
type Transaction struct {
	tx *sqlx.Tx
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	err := t.tx.Commit()
	t.tx = nil
	return err
}

// Rollback rolls the transaction back
func (t *Transaction) Rollback() error {
	err := t.tx.Rollback()
	t.tx = nil
	return err
}

// Active returns true if the transaction is active
func (t *Transaction) Active() bool {
	return t.tx == nil
}
