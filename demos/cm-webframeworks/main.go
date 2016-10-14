package main

import (
	"golang.org/x/net/context"
)

// main creates the database, creates some initial data, then starts the HTTP server
func main() {
	initDB()

	var frameworks []Framework
	if err := Frameworks.List(context.Background(), &frameworks); err != nil {
		panic(err)
	}

	if len(frameworks) == 0 {

		fx := Framework{
			Name:        "react",
			Description: "React, sometimes styled React.js or ReactJS, is an open-source JavaScript library for creating user interfaces that aims to address challenges encountered in developing single-page applications",
			URL:         "https://facebook.github.io/react/",
		}

		Frameworks.Insert(context.Background(), fx)

		fx = Framework{
			Name:        "riot",
			Description: "A React-like user interface micro-library",
			URL:         "http://riotjs.com/",
		}

		Frameworks.Insert(context.Background(), fx)

	}

	initHTTP()
}
