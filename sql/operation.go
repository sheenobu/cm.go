package sql

import (
	"fmt"

	"github.com/sheenobu/cm.go"
)

// UpdateOperation is the predicate which wraps a SQL update operation
type UpdateOperation struct {
	Column cm.ValueColumn
	Value  interface{}
}

// Apply modifies the collection to add the equal operation
func (op *UpdateOperation) Apply(c cm.Collection) error {
	col := c.(*Table)
	col.updateStatements = append(col.updateStatements,
		fmt.Sprintf("%s = ?", op.Column.Name()))

	col.updateValues = append(col.updateValues, op.Value)

	return nil
}
