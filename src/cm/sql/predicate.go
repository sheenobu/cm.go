package sql

import (
	"cm"
	"fmt"
)

// SqlEqPredicate is the predicate which wraps a SQL equal comparison
type SqlEqPredicate struct {
	Column cm.ValueColumn
	Value  interface{}
}

// begin Predicate implementation

// Apply modifies the collection to add the equal operation
func (pred *SqlEqPredicate) Apply(c cm.Collection) error {
	col := c.(*SqlTable)
	col.filterStatements = append(col.filterStatements,
		fmt.Sprintf("%s = ?", pred.Column.Name()))

	col.filterValues = append(col.filterValues, pred.Value)

	return nil
}

// end Predicate implementation

// SqlNotEqPredicate is the predicate which wraps a SQL not equal comparison
type SqlNotEqPredicate struct {
	Column cm.ValueColumn
	Value  interface{}
}

// begin Predicate implementation

// Apply modifies the collection to add the not equal operation
func (pred *SqlNotEqPredicate) Apply(c cm.Collection) error {
	col := c.(*SqlTable)
	col.filterStatements = append(col.filterStatements,
		fmt.Sprintf("%s != ?", pred.Column.Name()))

	col.filterValues = append(col.filterValues, pred.Value)

	return nil
}

// end Predicate implementation

// SqlLikePredicate is the predicate which wraps a SQL like comparison
type SqlLikePredicate struct {
	Column        cm.ValueColumn
	Value         interface{}
	CaseSensitive bool
}

// begin Predicate implementation

// Apply modifies the collection to add the like operation
func (pred *SqlLikePredicate) Apply(c cm.Collection) error {
	col := c.(*SqlTable)

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
