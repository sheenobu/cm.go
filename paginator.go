package cm

import "context"

// Paginator is the stateful object for allowing
// page based iteration of a collection.
type Paginator interface {

	// Apply pushes the current page of results into the given slice
	Apply(context.Context, interface{}) error

	// PageCount returns the number of pages
	PageCount() int

	// ItemCount returns the total number of items
	ItemCount() int

	// CurrentPage returns the current page
	CurrentPage() int

	// PerPageCount returns how many items are going to be on each page
	PerPageCount() int

	// Next moves to the next page, returning false if we are out of pages
	Next() bool

	// Prev moves to the previous page, returning false if we are out of pages
	Prev() bool
}
