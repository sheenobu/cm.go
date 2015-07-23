package sql

import (
	"cm"
)

// SqlValueColumn is a value column that contains metadata about the database column
type SqlValueColumn struct {
	name  string
	ctype string
}

// Column returns a column object given the name and type
func Column(name string, ctype string) cm.ValueColumn {
	return SqlValueColumn{name, ctype}
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

// end ValueColumn implementation
