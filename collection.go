package cm

import "context"

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

	// C returns the internal composite type, if it exists. Nil otherwise
	C() Collection
}
