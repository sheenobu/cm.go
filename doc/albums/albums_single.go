// +build ignore

package albums

import cm "github.com/sheenobu/cm.go"

var Albums AlbumsCollection

type AlbumsCollection struct {
	cm.Collection
	ID       cm.ValueColumn
	Artist   cm.ValueColumn
	Name     cm.ValueColumn
	Year     cm.ValueColumn
	Explicit cm.ValueColumn
}
