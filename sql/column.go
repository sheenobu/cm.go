package sql

import (
	"fmt"

	"github.com/sheenobu/cm.go"
)

// ValueColumn is a value column that contains metadata about the database column
type ValueColumn struct {
	name  string
	ctype string
	null  bool
	fns   map[string]func() interface{}
}

// Column returns a column object given the name and type
func Column(name string, ctype string) ValueColumn {
	return ValueColumn{name, ctype, true, make(map[string]func() interface{})}
}

// Varchar returns a column object given the name and size of string
func Varchar(name string, size int) ValueColumn {
	return Column(name,
		fmt.Sprintf("varchar(%d)", size))
}

// Integer returns an integer given the size
func Integer(name string, size int) ValueColumn {
	//TODO: make size a specific flag (BIGINT, SMALLINT, etc)
	return Column(name,
		fmt.Sprintf("integer"))
}

// PrimaryKey returns a value column that is a primary key
func (c ValueColumn) PrimaryKey() ValueColumn {
	c.ctype = c.ctype + " PRIMARY KEY "
	return c
}

// FromFunction returns a value column that is populated, on insert, from
// the given function
func (c ValueColumn) FromFunction(fn func() interface{}) ValueColumn {
	c.fns["insert"] = fn
	return c
}

// AutoIncrement returns the value column definition with autoincrement added
func (c ValueColumn) AutoIncrement() ValueColumn {
	c.ctype = c.ctype + " AUTOINCREMENT "
	return c
}

// NotNull returns a column object that does not allow null values
func (c ValueColumn) NotNull() ValueColumn {
	c.null = false
	return c
}

// Build builds the SQL column expression
func (c ValueColumn) Build() string {
	s := c.name + " " + c.ctype + " "
	if !c.null {
		s = s + " not null"
	}
	return s
}

// begin ValueColumn implementation

// Name returns the name of the SQL value column
func (c ValueColumn) Name() string {
	return c.name
}

// Type returns the type of the SQL value column
func (c ValueColumn) Type() string {
	return c.ctype
}

// Eq creates an equal predicate used for filtering
func (c ValueColumn) Eq(i interface{}) cm.Predicate {
	return &EqPredicate{c, i}
}

// NotEq creates a not equal predicate used for filtering
func (c ValueColumn) NotEq(i interface{}) cm.Predicate {
	return &NotEqPredicate{c, i}
}

// Like creates a like predicate used for filtering
func (c ValueColumn) Like(caseSensitive bool, i interface{}) cm.Predicate {
	return &LikePredicate{c, i, caseSensitive}
}

// Set creates a modify operation
func (c ValueColumn) Set(i interface{}) cm.Operation {
	return &UpdateOperation{c, i}
}

// end ValueColumn implementation
