package tx

import (
	"context"

	cm "github.com/sheenobu/cm.go"
)

type txkeytype int

var txkey txkeytype

// Supports returns true if the collection can support transactions
func Supports(col cm.Collection) bool {
	_, ok := col.C().(cm.Transactioner)
	return ok
}

// Active returns true if the context has an active transaction
func Active(ctx context.Context) bool {
	tx, ok := Get(ctx)
	if !ok || tx == nil {
		return false
	}

	return tx.Active()
}

// Begin attemps to start a transaction for the
// given collection
func Begin(ctx context.Context, col cm.Collection) (context.Context, error) {
	ctx, _, err := TryBegin(ctx, col)
	return ctx, err
}

// Rollback attemps to rollback the current transaction
func Rollback(ctx context.Context) error {
	t, ok := ctx.Value(txkey).(cm.Transaction)
	if !ok {
		panic("No transaction")
	}

	return t.Rollback()
}

// Commit attemps to commit the current transaction
func Commit(ctx context.Context) error {
	t, ok := ctx.Value(txkey).(cm.Transaction)
	if !ok {
		panic("No transaction")
	}

	return t.Commit()
}

// TryBegin attemps to start a transaction for the
// given collection, returning true if it succeeded
func TryBegin(ctx context.Context, col cm.Collection) (context.Context, bool, error) {
	txer, ok := col.C().(cm.Transactioner)
	if !ok {
		return ctx, false, nil
	}

	t, err := txer.BeginTx()
	if err != nil {
		return ctx, false, err
	}

	ctx = context.WithValue(ctx, txkey, t)

	return ctx, true, err
}

// Get tries to get the transaction from the context
func Get(ctx context.Context) (cm.Transaction, bool) {
	t, ok := ctx.Value(txkey).(cm.Transaction)
	return t, ok
}
