// +build ignore

// Querying for all released in the year 2011
var albumList []albums.Album
err = Albums.Filter(Albums.Year.Eq(2011)).List(ctx, &albumList)

