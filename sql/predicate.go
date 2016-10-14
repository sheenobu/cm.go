package sql

import (
	"fmt"

	"github.com/sheenobu/cm.go"
)

// EqPredicate is the predicate which wraps a SQL equal comparison
type EqPredicate struct {
	Column cm.ValueColumn
	Value  interface{}
}

// begin Predicate implementation

// Apply modifies the collection to add the equal operation
func (pred *EqPredicate) Apply(c cm.Collection) error {
	col := c.(*Table)
	col.filterStatements = append(col.filterStatements,
		fmt.Sprintf("%s = ?", pred.Column.Name()))

	col.filterValues = append(col.filterValues, pred.Value)

	return nil
}

// end Predicate implementation

// NotEqPredicate is the predicate which wraps a SQL not equal comparison
type NotEqPredicate struct {
	Column cm.ValueColumn
	Value  interface{}
}

// begin Predicate implementation

// Apply modifies the collection to add the not equal operation
func (pred *NotEqPredicate) Apply(c cm.Collection) error {
	col := c.(*Table)
	col.filterStatements = append(col.filterStatements,
		fmt.Sprintf("%s != ?", pred.Column.Name()))

	col.filterValues = append(col.filterValues, pred.Value)

	return nil
}

// end Predicate implementation

// LikePredicate is the predicate which wraps a SQL like comparison
type LikePredicate struct {
	Column        cm.ValueColumn
	Value         interface{}
	CaseSensitive bool
}

// begin Predicate implementation

// Apply modifies the collection to add the like operation
func (pred *LikePredicate) Apply(c cm.Collection) error {
	col := c.(*Table)

	like := "like"

	if pred.CaseSensitive {
		like = "ilike"
	}

	col.filterStatements = append(col.filterStatements,
		fmt.Sprintf("%s %s ?", pred.Column.Name(), like))

	col.filterValues = append(col.filterValues, pred.Value)

	return nil
}

// end Predicate implementation
