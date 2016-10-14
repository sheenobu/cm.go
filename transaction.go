package cm

// Transactioner marks a collection as something that
// can support starting of transactions.
type Transactioner interface {
	BeginTx() (Transaction, error)
}

// A Transaction is a unit of work against a database
type Transaction interface {
	Commit() error
	Rollback() error
	Active() bool
}
