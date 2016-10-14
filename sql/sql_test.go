package sql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sheenobu/cm.go"

	"os"
	"testing"

	"golang.org/x/net/context"
)

// Car defines the model object we are using.
type Car struct {
	ID      *int
	Slug    *string
	CarMake string
	Model   string
	Year    int64
}

// _Cars is the collection for the model.
type _Cars struct {
	cm.Collection
	ID      cm.ValueColumn
	Slug    cm.ValueColumn
	CarMake cm.ValueColumn
	Model   cm.ValueColumn
	Year    cm.ValueColumn
}

func constant() interface{} {
	return "hello"
}

// createCars resets the database and initializes the cars structure
func createCars() *_Cars {
	os.Remove("cars.db")
	db, _ := sqlx.Connect("sqlite3", "cars.db")
	return &_Cars{
		Collection: New(db, "CARS"),
		ID:         Integer("id", 10).PrimaryKey().AutoIncrement(),
		Slug:       Varchar("slug", 32).FromFunction(constant),
		CarMake:    Varchar("make", 100).NotNull(),
		Model:      Varchar("model", 100).NotNull(),
		Year:       Column("year", "number").NotNull(),
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

	if s := cars[0].ID; *s != 1 {
		t.Errorf("Expected cars.Id to be 1, got %v", *s)
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

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, slug, model, make, year)
		values ('1', '', 'Honda', 'Civic', '1993'),
	           ('2', '', 'Toyota', 'Corolla', '1993');
	`)

	if err != nil {
		t.Error(err)
	}

	var cars []Car

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

	var cars []Car

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

	var cars []Car

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

	var cars []Car

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

// TestUpdateFilterL tests the update function on filtered items
func TestUpdateFilterL(t *testing.T) {

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

	var cars []Car

	err = Cars.Filter(Cars.Model.Eq("Honda")).Edit(Cars.Model.Set("Whatever")).Update(ctx)

	if err != nil {
		t.Error(err)
	}

	cars = make([]Car, 0)

	if err = Cars.Filter(Cars.Model.Eq("Honda")).List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 0 {
		t.Errorf("Expected Honda cars to be of length 0, got %d", i)
	}

	cars = make([]Car, 0)

	if err = Cars.Filter(Cars.Model.Eq("Whatever")).List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected Whatever cars to be of length 1, got %d", i)
	}

}

// TestUpdateFilterR tests the update function on filtered items
func TestUpdateFilterR(t *testing.T) {

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

	var cars []Car

	err = Cars.Edit(Cars.Model.Set("Whatever")).Filter(Cars.Model.Eq("Honda")).Update(ctx)

	if err != nil {
		t.Error(err)
	}

	cars = make([]Car, 0)

	if err = Cars.Filter(Cars.Model.Eq("Honda")).List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 0 {
		t.Errorf("Expected Honda cars to be of length 0, got %d", i)
	}

	cars = make([]Car, 0)

	if err = Cars.Filter(Cars.Model.Eq("Whatever")).List(ctx, &cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 1 {
		t.Errorf("Expected Whatever cars to be of length 1, got %d", i)
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

	var cars []Car

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

	var cars []Car

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

// TestInsert tests the inesrt function
func TestInsert(t *testing.T) {
	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
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
}

// TestPagination tests the pagination function
func TestPaginationDevideOddly(t *testing.T) {
	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Accord Gen1', '1976'),
		('2', 'Honda', 'Accord Gen1', '1977'),
		('3', 'Honda', 'Accord Gen1', '1978'),
		('4', 'Honda', 'Accord Gen1', '1979'),
		('5', 'Honda', 'Accord Gen1', '1980'),
		('6', 'Honda', 'Accord Gen2', '1981'),
		('7', 'Honda', 'Accord Gen3', '1982'),
		('8', 'Honda', 'Accord Gen3', '1983'),
		('9', 'Honda', 'Accord Gen3', '1984'),
		('10', 'Honda', 'Accord Gen4', '1985'),
		('11', 'Honda', 'Accord Gen4', '1986'),
		('12', 'Honda', 'Accord Gen5', '1987'),
		('13', 'Honda', 'Accord Gen5', '1988'),
		('14', 'Honda', 'Accord Gen5', '1989');
	`)

	if err != nil {
		t.Error(err)
	}

	page, err := Cars.Page(ctx, 3)

	if err != nil {
		t.Error(err)
	}

	if page == nil {
		t.Error("Page is empty")
	}

	if i := page.PageCount(); i != 5 {
		t.Errorf("Expected page count to be 5, was %d", i)
	}

	if i := page.CurrentPage(); i != 0 {
		t.Errorf("Expected initial page to be 0, was %d", i)
	}

	if ok := page.Prev(); ok {
		t.Error("Page Previous should have returned false")
	}

	if ok := page.Next(); !ok {
		t.Errorf("Page next should have succeeded")
	}

	if i := page.CurrentPage(); i != 1 {
		t.Errorf("Expected current page after Next to be 1, was %d", i)
	}

	var cars []Car

	if err := page.Apply(&cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 3 {
		t.Errorf("Expected car page length to be 3, was %d", i)
	} else {
		c0 := cars[0]
		c1 := cars[1]
		c2 := cars[2]

		if c0.Year != 1979 {
			t.Errorf("Expected car year to be 1979, was %d", c0.Year)
		}
		if c1.Year != 1980 {
			t.Errorf("Expected car year to be 1980, was %d", c1.Year)
		}
		if c2.Year != 1981 {
			t.Errorf("Expected car year to be 1981, was %d", c2.Year)
		}
	}

	ok := true

	var i = 0

	// reset the pagination
	for ok && i < 100 {
		ok = page.Prev()
		i++
	}
	if i >= 100 {
		t.Errorf("page.Prev() continuously ran")
	}

	ok = true

	var lengths []int

	i = 0
	for ok && i < 100 {
		cars = make([]Car, 0)
		page.Apply(&cars)
		lengths = append(lengths, len(cars))
		ok = page.Next()
		i++
	}
	if i >= 100 {
		t.Errorf("page.Prev() continuously ran")
	}

	if i := len(lengths); i != 5 {
		t.Errorf("Should have iterated over 5 pages, but iterated over %d", i)
	}

	if len(lengths) == 5 {
		if lengths[0] != 3 {
			t.Errorf("First page should have had 3 items")
		}
		if lengths[1] != 3 {
			t.Errorf("Second page should have had 3 items")
		}
		if lengths[2] != 3 {
			t.Errorf("Third page should have had 3 items")
		}
		if lengths[3] != 3 {
			t.Errorf("Fourth page should have had 3 items, has %d", lengths[3])
		}
		if lengths[4] != 2 {
			t.Errorf("Fifth page should have had 2 items, has %d", lengths[4])
		}
	}

}

// TestPaginationDivideEvenly tests the pagination function
func TestPaginationDivideEvenly(t *testing.T) {
	ctx := context.Background()

	Cars := createCars()
	err := Cars.Init(Cars)

	if err != nil {
		t.Error(err)
	}

	err = Cars.ExecRaw(context.Background(), `insert into CARS (id, model, make, year)
		values ('1', 'Honda', 'Accord Gen1', '1976'),
		('2', 'Honda', 'Accord Gen1', '1977'),
		('3', 'Honda', 'Accord Gen1', '1978'),
		('4', 'Honda', 'Accord Gen1', '1979'),
		('5', 'Honda', 'Accord Gen1', '1980'),
		('6', 'Honda', 'Accord Gen2', '1981'),
		('7', 'Honda', 'Accord Gen3', '1982'),
		('8', 'Honda', 'Accord Gen3', '1983'),
		('9', 'Honda', 'Accord Gen3', '1984'),
		('10', 'Honda', 'Accord Gen4', '1985'),
		('11', 'Honda', 'Accord Gen4', '1986'),
		('12', 'Honda', 'Accord Gen5', '1987'),
		('13', 'Honda', 'Accord Gen5', '1988'),
		('14', 'Honda', 'Accord Gen5', '1989')
	`)

	if err != nil {
		t.Error(err)
	}

	page, err := Cars.Page(ctx, 2)

	if err != nil {
		t.Error(err)
	}

	if page == nil {
		t.Error("Page is empty")
	}

	if i := page.PageCount(); i != 7 {
		t.Errorf("Expected page count to be 7, was %d", i)
	}

	if i := page.CurrentPage(); i != 0 {
		t.Errorf("Expected initial page to be 0, was %d", i)
	}

	if ok := page.Prev(); ok {
		t.Error("Page Previous should have returned false")
	}

	if ok := page.Next(); !ok {
		t.Errorf("Page next should have succeeded")
	}

	if i := page.CurrentPage(); i != 1 {
		t.Errorf("Expected current page after Next to be 1, was %d", i)
	}

	var cars []Car

	if err := page.Apply(&cars); err != nil {
		t.Error(err)
	}

	if i := len(cars); i != 2 {
		t.Errorf("Expected car page length to be 2, was %d", i)
	} else {
		c0 := cars[0]
		c1 := cars[1]

		if c0.Year != 1978 {
			t.Errorf("Expected car year to be 1979, was %d", c0.Year)
		}
		if c1.Year != 1979 {
			t.Errorf("Expected car year to be 1980, was %d", c1.Year)
		}
	}

	ok := true

	var i = 0
	// reset the pagination
	for ok && i < 100 {
		ok = page.Prev()
		i++
	}
	if i >= 100 {
		t.Errorf("page.prev never stopped looping")
	}

	if page.CurrentPage() != 0 {
		t.Errorf("Current page should have been 0, was %d\n", page.CurrentPage())
	}

	ok = true

	var lengths []int

	i = 0
	for ok && i < 100 {
		cars = make([]Car, 0)
		page.Apply(&cars)
		lengths = append(lengths, len(cars))
		ok = page.Next()
		i++
	}
	if i >= 100 {
		t.Errorf("page.prev never stopped looping")
	}

	if i := len(lengths); i != 7 {
		t.Errorf("Should have iterated over 7  pages, but iterated over %d", i)
	}

	for idx, i := range lengths {
		if i != 2 {
			t.Errorf("Page %d should have had 2 items, had %d", idx, i)
		}
	}

}
