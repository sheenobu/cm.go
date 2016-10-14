package sql

import (
	"context"
	"testing"

	"github.com/sheenobu/cm.go/tx"
)

// TestInsertWithRollback tests the insert function
// then rolls back the transaction
func TestInsertWithRollback(t *testing.T) {
	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	if !tx.Supports(Cars) {
		t.Errorf("The collection does not support transactions.")
		return
	}

	ctx, err = tx.Begin(ctx, Cars)
	if err != nil {
		t.Errorf("Error starting transaction: '%v'", err)
	}

	car := &Car{
		CarMake: "Honda",
		Model:   "Civic",
		Year:    1997,
	}

	err = Cars.Insert(ctx, car)
	if err != nil {
		t.Error(err)
	}

	var cars []Car
	if err = Cars.List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected all cars to be of length 1, got %d", i)
	}

	if len(cars) == 1 {
		if cars[0].Model != "Civic" {
			t.Errorf("Expected remaining car to be Civic, was %s", cars[0].Model)
		}
		if cars[0].ID == nil {
			t.Errorf("No ID on car")
		}
		if cars[0].Slug != nil {
			if *(cars[0].Slug) != "hello" {
				t.Errorf("Expected remaining car slug to be 'hello', was %s", *cars[0].Slug)
			}
		} else {
			t.Errorf("Expected slug to be non-nil, was nil")
		}
	}

	err = tx.Rollback(ctx)
	if err != nil {
		t.Errorf("Error rolling back transaction: '%v'", err)
	}

	if err = Cars.List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected all cars to be of length 0 after transaction rollback, got %d", i)
	}

}
