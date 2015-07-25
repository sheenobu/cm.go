package main

// Framework is our simple model object
type Framework struct {
	Id          string // Id is an integer but defined as a string (for internal buggy reasons within the sql support)
	Name        string
	Description string
	Url         string
}
