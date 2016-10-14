package cm

// Operation defines a piece of code that
// modifies a writing statement
type Operation interface {
	Apply(c Collection) error
}
