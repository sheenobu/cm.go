package sql

import (
	"fmt"

	"github.com/sheenobu/cm.go"
)

// SqlValueColumn is a value column that contains metadata about the database column
type SqlValueColumn struct {
	name  string
	ctype string
	null  bool
}

// Column returns a column object given the name and type
func Column(name string, ctype string) SqlValueColumn {
	return SqlValueColumn{name, ctype, true}
}

// Varchar returns a column object given the name and size of string
func Varchar(name string, size int) SqlValueColumn {
	return Column(name,
		fmt.Sprintf("varchar(%d)", size))
}

// Integer returns an integer given the size
func Integer(name string, size int) SqlValueColumn {
	//TODO: make size a specific flag (BIGINT, SMALLINT, etc)
	return Column(name,
		fmt.Sprintf("integer"))
}

func (c SqlValueColumn) PrimaryKey() SqlValueColumn {
	c.ctype = c.ctype + " PRIMARY KEY "
	return c
}

func (c SqlValueColumn) AutoIncrement() SqlValueColumn {
	c.ctype = c.ctype + " AUTOINCREMENT "
	return c
}

// NotNull returns a column object that does not allow null values
func (c SqlValueColumn) NotNull() SqlValueColumn {
	c.null = false
	return c
}

// Build builds the SQL column expression
func (c SqlValueColumn) Build() string {
	s := c.name + " " + c.ctype + " "
	if !c.null {
		s = s + " not null"
	}
	return s
}

// begin ValueColumn implementation

// Name returns the name of the SQL value column
func (s SqlValueColumn) Name() string {
	return s.name
}

// Type returns the type of the SQL value column
func (s SqlValueColumn) Type() string {
	return s.ctype
}

// Eq creates an equal predicate used for filtering
func (s SqlValueColumn) Eq(i interface{}) cm.Predicate {
	return &SqlEqPredicate{s, i}
}

// NotEq creates a not equal predicate used for filtering
func (s SqlValueColumn) NotEq(i interface{}) cm.Predicate {
	return &SqlNotEqPredicate{s, i}
}

// Like creates a like predicate used for filtering
func (s SqlValueColumn) Like(caseSensitive bool, i interface{}) cm.Predicate {
	return &SqlLikePredicate{s, i, caseSensitive}
}

// Set creates a modify operation
func (s SqlValueColumn) Set(i interface{}) cm.Operation {
	return &SqlUpdateOperation{s, i}
}

// end ValueColumn implementation
