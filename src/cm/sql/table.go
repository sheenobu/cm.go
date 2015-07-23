package sql

import (
	"fmt"
	"github.com/fatih/camelcase"
	"github.com/jmoiron/sqlx"
	"strings"

	"golang.org/x/net/context"
	"reflect"

	"cm"
)

// SqlTable defines the Collection that interacts with a sqlx.DB connection
type SqlTable struct {
	Table string

	database         *sqlx.DB
	columns          []string
	filterStatements []string
	values           []interface{}
}

// New defines a new SqlTable based on the database connection and table name.
func New(db *sqlx.DB, table string) *SqlTable {
	return &SqlTable{
		Table:            table,
		database:         db,
		columns:          make([]string, 0),
		filterStatements: make([]string, 0),
		values:           make([]interface{}, 0),
	}
}

// begin Collection implementation

// Init initializes the sql table by performing reflection operations
// on the interface
func (sql *SqlTable) Init(iface interface{}) error {
	columns := make([]string, 0)

	s := reflect.ValueOf(iface).Elem()
	st := reflect.TypeOf(iface).Elem()

	for i := 0; i < st.NumField(); i++ {
		f := s.Field(i)
		fx := st.Field(i)

		if f.Type().Name() == "ValueColumn" {

			vc := f.Interface().(cm.ValueColumn)

			columns = append(columns, vc.Name()+" "+vc.Type())

			cs := camelcase.Split(fx.Name)

			for i := 0; i != len(cs); i++ {
				cs[i] = strings.ToLower(cs[i])
			}

			sql.columns = append(sql.columns,
				vc.Name()+" as "+strings.ToLower(fx.Name))

		}
	}

	createTable := fmt.Sprintf(`
		create table %s
			(%s);
	`, sql.Table, strings.Join(columns, ", "))

	_, err := sql.database.Exec(createTable)

	return err
}

// Filter filters the collection on the given predicate
func (sql SqlTable) Filter(pred cm.Predicate) cm.Collection {
	pred.Apply(&sql)
	return &sql
}

// List resolves the collection and applies the results to the given slice pointer
func (sql SqlTable) List(ctx context.Context, list interface{}) (err error) {

	selectQ := "select " + strings.Join(sql.columns, ", ") + " from " + sql.Table

	if len(sql.filterStatements) == 0 {
		err = sql.database.Select(list, selectQ)
	} else {
		q := selectQ + " where " + strings.Join(sql.filterStatements, " and ")
		err = sql.database.Select(list, q, sql.values...)
	}

	return err
}

// Single resolves the collection and applies the results to the given slice pointer.
func (sql SqlTable) Single(ctx context.Context, single interface{}) (err error) {
	selectQ := "select " + strings.Join(sql.columns, ", ") + " from " + sql.Table

	if len(sql.filterStatements) == 0 {
		err = sql.database.Select(single, selectQ+" limit 1")
	} else {
		err = sql.database.Select(single, selectQ+" where "+strings.Join(sql.filterStatements, " and ")+" limit 1", sql.values...)
	}
	return err
}

// Delete resolves the collection using the filters and deletes the filtered elements.
func (sql SqlTable) Delete(ctx context.Context) (err error) {
	deleteQ := "delete from " + sql.Table

	if len(sql.filterStatements) == 0 {
		_, err = sql.database.Exec(deleteQ)
	} else {
		_, err = sql.database.Exec(deleteQ+" where "+strings.Join(sql.filterStatements, " and "), sql.values...)
	}

	return err
}

// ExecRaw executes a raw statement against the collection
func (sql SqlTable) ExecRaw(ctx context.Context, query string) error {
	_, err := sql.database.Exec(query)
	return err
}

// end Collection implementation
