package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sheenobu/cm.go/tx"
)

// main creates the database, creates some initial data, then starts the HTTP server
func main() {
	initDB()

	err := run()
	if err != nil {
		fmt.Printf("Error: '%v'\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() (err error) {
	var ctx = context.Background()

	defer func() {
		if tx.Active(ctx) {
			if err != nil {
				tx.Rollback(ctx)
			} else {
				tx.Commit(ctx)
			}
		}
	}()

	ctx, err = tx.Begin(ctx, Frameworks)
	if err != nil {
		return
	}

	var frameworks []Framework
	if err = Frameworks.List(ctx, &frameworks); err != nil {
		return
	}

	if len(frameworks) == 0 {

		fx := Framework{
			Name:        "react",
			Description: "React, sometimes styled React.js or ReactJS, is an open-source JavaScript library for creating user interfaces that aims to address challenges encountered in developing single-page applications",
			URL:         "https://facebook.github.io/react/",
		}

		Frameworks.Insert(ctx, fx)

		fx = Framework{
			Name:        "riot",
			Description: "A React-like user interface micro-library",
			URL:         "http://riotjs.com/",
		}

		Frameworks.Insert(ctx, fx)

	}

	if err = tx.Commit(ctx); err != nil {
		return
	}

	err = initHTTP()

	return
}
