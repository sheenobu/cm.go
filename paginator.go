package cm

import "context"

// Paginator defines an object which allows pagination, iteration
// on a collection
type Paginator interface {
	Apply(context.Context, interface{}) error

	PageCount() int
	ItemCount() int

	CurrentPage() int
	PerPageCount() int

	Next() bool
	Prev() bool
}
