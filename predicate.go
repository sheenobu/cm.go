package cm

// Predicate defines a piece of code that
// modifies a reading statement.
type Predicate interface {
	Apply(c Collection) error
}
