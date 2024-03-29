// +build ignore

package albums

import cm "github.com/sheenobu/cm.go"

// Album is the model for music album database
type Album struct {
	ID       *int
	Artist   string
	Name     string
	Explicit bool
	Year     int64
}

// AlbumsCollection defines the columns and operations for
// the Album model
type AlbumsCollection struct {
	cm.Collection
	ID       cm.ValueColumn
	Artist   cm.ValueColumn
	Name     cm.ValueColumn
	Year     cm.ValueColumn
	Explicit cm.ValueColumn
}
