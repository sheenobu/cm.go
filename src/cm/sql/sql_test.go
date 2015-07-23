package sql

import (
	"cm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/net/context"
	"os"
	"testing"
)

// Car defines the model object we are using.
type Car struct {
	PrimaryKey string
	CarMake    string
	Model      string
	Year       int64
}

// _Cars is the collection for the model.
type _Cars struct {
	cm.Collection
	PrimaryKey cm.ValueColumn
	CarMake    cm.ValueColumn
	Model      cm.ValueColumn
	Year       cm.ValueColumn
}

// createCars resets the database and initializes the cars structure
func createCars() *_Cars {
	os.Remove("cars.db")
	db, _ := sqlx.Connect("sqlite3", "cars.db")
	return &_Cars{
		Collection: New(db, "CARS"),
		PrimaryKey: Column("id", "varchar(33) not null primary key"),
		CarMake:    Column("make", "varchar(100) not null"),
		Model:      Column("model", "varchar(100) not null"),
		Year:       Column("year", "number not null"),
	}
}

// TestInit initializes the database
func TestInit(t *testing.T) {
	Cars := createCars()
	err := Cars.Init(Cars)
	if err != nil {
		t.Error(err)
	}
}

// TestAllEmpty ensures the collection is empty when initially created.
func TestAllEmpty(t *testing.T) {
	Cars := createCars()
	err := Cars.Init(Cars)
	if err != nil {
		t.Error(err)
	}

	var cars []Car

	err = Cars.List(context.Background(), &cars)
	if err != nil {
		t.Error(err)
	}

	if len(cars) > 0 {
		t.Errorf("Car collection should be empty")
	}

	err = Cars.List(context.Background(), &cars)
	if err != nil {
		t.Error(err)
	}

	if len(cars) > 0 {
		t.Errorf("Car collection should be empty")
	}

}

// TestAll tests the All function.
func TestAll(t *testing.T) {
	Cars := createCars()
	err := Cars.Init(Cars)
	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Civic', '1993'),
	           ('2', 'Toyota', 'Corolla', '1993');
	`)

	if err != nil {
		t.Error(err)
	}

	var cars []Car

	err = Cars.List(context.Background(), &cars)
	if err != nil {
		t.Error(err)
		return
	}

	if i := len(cars); i != 2 {
		t.Errorf("Expected cars to be of length 2, got %d", i)
		return
	}

	if s := cars[0].PrimaryKey; s != "1" {
		t.Errorf("Expected cars.Id to be 1, got %s", s)
	}
}

// TestEqFilter tests the filter function with equality
func TestEqFilter(t *testing.T) {

	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Civic', '1993'),
	           ('2', 'Toyota', 'Corolla', '1993');
	`)

	if err != nil {
		t.Error(err)
	}

	cars := make([]Car, 0)

	err = Cars.Filter(Cars.Year.Eq("1993")).List(ctx, &cars)

	if err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 2 {
		t.Errorf("Expected cars made in 1993 to be of length 2, got %d", i)
	}

	cars = make([]Car, 0)

	err = Cars.Filter(Cars.Model.Eq("Honda")).List(ctx, &cars)

	if err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected Honda model cars to be of length 1, got %d", i)
	}

	cars = make([]Car, 0)

	// This is an AND, basically
	err = Cars.Filter(Cars.Year.Eq("1993")).Filter(Cars.Model.Eq("Honda")).List(ctx, &cars)

	if err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected Honda model cars to be of length 1, got %d", i)
	}

}

// TestNotEqFilter tests the filter function with not equality
func TestNotEqFilter(t *testing.T) {

	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Civic', '1993'),
	           ('2', 'Toyota', 'Corolla', '1993');
	`)

	if err != nil {
		t.Error(err)
	}

	cars := make([]Car, 0)

	err = Cars.Filter(Cars.Year.NotEq("1993")).List(ctx, &cars)

	if err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 0 {
		t.Errorf("Expected cars not made in 1993 to be of length 0, got %d", i)
	}

	cars = make([]Car, 0)

	err = Cars.Filter(Cars.Model.NotEq("Honda")).List(ctx, &cars)

	if err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected Non-Honda model cars to be of length 1, got %d", i)
	}

	if len(cars) > 0 {
		if s := cars[0].Model; s != "Toyota" {
			t.Errorf("Expected Non-Honda model car to be Toyota, got %s", s)
		}
	}

	cars = make([]Car, 0)

	// This is an AND, basically
	err = Cars.Filter(Cars.Year.Eq("1993")).Filter(Cars.Model.NotEq("Honda")).List(ctx, &cars)

	if err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected Non honda model cars to be of length 1, got %d", i)
	}

	if len(cars) > 0 {
		if s := cars[0].Model; s != "Toyota" {
			t.Errorf("Expected Non-Honda model car to be Toyota, got %s", s)
		}
	}
}

// TestLikeFilter tests the filter function with like wildcard equality
func TestLikeFilter(t *testing.T) {

	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Civic', '1993'),
	           ('2', 'Toyota', 'Corolla', '1993');
	`)

	if err != nil {
		t.Error(err)
	}

	cars := make([]Car, 0)

	err = Cars.Filter(Cars.Model.Like(false, "%a")).List(ctx, &cars)

	if err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 2 {
		t.Errorf("Expected all cars to be of length 2, got %d", i)
	}

	cars = make([]Car, 0)

	// This is an AND, basically
	err = Cars.Filter(Cars.Year.Eq("1993")).Filter(Cars.Model.Like(false, "H%")).List(ctx, &cars)

	if err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected honda model cars to be of length 1, got %d", i)
	}

	if len(cars) > 0 {
		if s := cars[0].Model; s != "Honda" {
			t.Errorf("Expected Honda model car, got %s", s)
		}
	}
}

// TestUpdateAll tests the update function on all items
func TestUpdateAll(t *testing.T) {

	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Civic', '1993'),
	           ('2', 'Toyota', 'Corolla', '1993');
	`)

	if err != nil {
		t.Error(err)
	}

	cars := make([]Car, 0)

	err = Cars.Edit(Cars.Model.Set("Whatever")).Update(ctx)

	if err != nil {
		t.Error(err)
	}

	if err = Cars.List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 2 {
		t.Errorf("Expected all cars to be of length 2, got %d", i)
	}

	if len(cars) == 2 {
		for _, car := range cars {
			if car.Model != "Whatever" {
				t.Errorf("Expected all cars to have model of Whatever")
			}
		}
	}
}

// TestDeleteAll tests the delete function on all items
func TestDeleteAll(t *testing.T) {

	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Civic', '1993'),
	           ('2', 'Toyota', 'Corolla', '1993');
	`)

	if err != nil {
		t.Error(err)
	}

	cars := make([]Car, 0)

	err = Cars.Delete(ctx)

	if err != nil {
		t.Error(err)
	}

	if err = Cars.List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 0 {
		t.Errorf("Expected all cars to be of length 0, got %d", i)
	}
}

// TestDeleteFilter tests the delete function alongside a filter
func TestDeleteFilter(t *testing.T) {
	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Civic', '1993'),
	           ('2', 'Toyota', 'Corolla', '1993');
	`)

	if err != nil {
		t.Error(err)
	}

	cars := make([]Car, 0)

	err = Cars.Filter(Cars.Model.Eq("Honda")).Delete(ctx)

	if err != nil {
		t.Error(err)
	}

	if err = Cars.List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected all cars to be of length 1, got %d", i)
	}

	if len(cars) == 1 {
		if cars[0].Model != "Toyota" {
			t.Errorf("Expected remaining car to be Toyota, was %s", cars[0].Model)
		}
	}
}
