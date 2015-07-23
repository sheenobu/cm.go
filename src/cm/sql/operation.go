package sql

import (
	"cm"
	"fmt"
)

// SqlUpdateOperation is the predicate which wraps a SQL update operation
type SqlUpdateOperation struct {
	Column cm.ValueColumn
	Value  interface{}
}

// Apply modifies the collection to add the equal operation
func (op *SqlUpdateOperation) Apply(c cm.Collection) error {
	col := c.(*SqlTable)
	col.updateStatements = append(col.updateStatements,
		fmt.Sprintf("%s = ?", op.Column.Name()))

	col.updateValues = append(col.updateValues, op.Value)

	return nil
}
