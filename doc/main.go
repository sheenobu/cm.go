// +build ignore

package main

import (
	"context"
	"github.com/sheenobu/cm.go/doc/albums"
)

// Albums is an optional shortcut for our example code
var Albums = albums.Collection

func main() {
	// for potential transactions
	ctx := context.Background()

	// insertion
	err = Albums.Insert(ctx, &albums.Album{
		Artist:   "Childish Gambino",
		Name:     "Camp",
		Year:     2011,
		Explicit: true,
	})

	err = Albums.Insert(ctx, &albums.Album{
		Artist:   "Sleater-Kinney",
		Name:     "No Cities To Love",
		Year:     2015,
		Explicit: false,
	})

	// Querying for all released in the year 2011
	var albumList []albums.Album
	err = Albums.Filter(Albums.Year.Eq(2011)).List(ctx, &albumList)

	// Updating all explicit albums to non-explicit
	filter := Albums.Explicit.Eq(true)
	setter := Albums.Explicit.Set(false)
	err = Albums.Filter(filter).Edit(setter).Update(ctx)

	// Deletion by ID
	err = Albums.Filter(Albums.ID.Eq(0)).Delete(ctx)
}
