package cm

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
