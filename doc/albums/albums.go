// +build ignore

package albums

import (
	"github.com/sheenobu/cm.go"
)

// Album defines the model for the music album
type Album struct {
	ID       *int
	Artist   string
	Name     string
	Explicit bool
	Year     int64
}

// _Albums is the collection for the model.
type _Albums struct {
	cm.Collection
	ID       cm.ValueColumn
	Artist   cm.ValueColumn
	Name     cm.ValueColumn
	Year     cm.ValueColumn
	Explicit cm.ValueColumn
}
