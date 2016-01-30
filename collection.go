package cm

import (
	"golang.org/x/net/context"
)

// ValueColumn defines a column defined on a collection.
type ValueColumn interface {

	// Name is the column name
	Name() string

	// Type is the column type (sql type, redis type, etc)
	Type() string

	// Eq returns a predicate that, when applied to a Collection, filters out matching elements
	Eq(interface{}) Predicate

	// NotEq returns a predicate that, when applied to a Collection, filters out non-matching elements
	NotEq(interface{}) Predicate

	// Like returns a predicate that, when applied to a Collection, filters out wildcard matching elements
	Like(bool, interface{}) Predicate

	// Set returns the operation that, when applied to a collection, updates the elements
	Set(interface{}) Operation
}

// Predicate defines a piece of code that
// modifies a reading statement.
type Predicate interface {
	Apply(c Collection) error
}

// Operation defines a piece of code that
// modifies a writing statement
type Operation interface {
	Apply(c Collection) error
}

// Paginator defines an object which allows pagination, iteration
// on a collection
type Paginator interface {
	Apply(interface{}) error

	PageCount() int
	ItemCount() int

	CurrentPage() int
	PerPageCount() int

	Next() bool
	Prev() bool
}

// Collection defines a grouping of related entities that can be
// operated on.
type Collection interface {
	Init(interface{}) error

	ExecRaw(context.Context, string) error

	// Modify Operations
	Filter(Predicate) Collection
	Edit(Operation) Collection

	/*

		// Entity Operations
		EntityInsert(context.Context, *interface{}) error
		EntityUpdate(context.Context, *interface{}) error
		EntityRemove(context.Context, *interface{}) error

		// Batch Operations
	*/

	Insert(context.Context, interface{}) error
	Delete(context.Context) error
	Update(context.Context) error

	// Read Operations
	List(context.Context, interface{}) error
	Page(context.Context, int) (Paginator, error)
	Single(context.Context, interface{}) error
}