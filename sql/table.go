package sql

import (
	"fmt"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/jmoiron/sqlx"
	"github.com/sheenobu/cm.go"

	"reflect"

	"golang.org/x/net/context"
)

// SqlTable defines the Collection that interacts with a sqlx.DB connection
type SqlTable struct {
	Table string

	database     *sqlx.DB
	asColumns    []string
	insColumns   []string
	namedColumns []string

	filterStatements []string
	updateStatements []string
	updateValues     []interface{}
	filterValues     []interface{}
}

// New defines a new SqlTable based on the database connection and table name.
func New(db *sqlx.DB, table string) *SqlTable {
	return &SqlTable{
		Table:    table,
		database: db,

		asColumns:    make([]string, 0),
		insColumns:   make([]string, 0),
		namedColumns: make([]string, 0),

		filterStatements: make([]string, 0),
		updateStatements: make([]string, 0),
		updateValues:     make([]interface{}, 0),
		filterValues:     make([]interface{}, 0),
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

			vc := f.Interface().(SqlValueColumn)

			columns = append(columns, vc.Build())

			cs := camelcase.Split(fx.Name)

			for i := 0; i != len(cs); i++ {
				cs[i] = strings.ToLower(cs[i])
			}

			sql.asColumns = append(sql.asColumns,
				vc.Name()+" as "+strings.ToLower(fx.Name))

			sql.insColumns = append(sql.insColumns,
				vc.Name())

			sql.namedColumns = append(sql.namedColumns,
				":"+strings.ToLower(fx.Name))

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

// Edit updates the collection for the given operation
func (sql SqlTable) Edit(op cm.Operation) cm.Collection {
	op.Apply(&sql)
	return &sql
}

// List resolves the collection and applies the results to the given slice pointer
func (sql SqlTable) List(ctx context.Context, list interface{}) (err error) {

	selectQ := "select " + strings.Join(sql.asColumns, ", ") + " from " + sql.Table

	if len(sql.filterStatements) == 0 {
		err = sql.database.Select(list, selectQ)
	} else {
		q := selectQ + " where " + strings.Join(sql.filterStatements, " and ")
		err = sql.database.Select(list, q, sql.filterValues...)
	}

	return err
}

// Page returns a pagination object which is responsible for resolving the collection
func (sql SqlTable) Page(ctx context.Context, perPageCount int) (cm.Paginator, error) {

	countQ := "select count(" + sql.insColumns[0] + ") from " + sql.Table
	selectQ := "select " + strings.Join(sql.asColumns, ", ") + " from " + sql.Table

	if len(sql.filterStatements) != 0 {
		selectQ = selectQ + " where " + strings.Join(sql.filterStatements, " and ")
		countQ = countQ + " where " + strings.Join(sql.filterStatements, " and ")
	}

	sp := SqlPaginator{
		database:     sql.database,
		selectQuery:  selectQ,
		countQuery:   countQ,
		perPageCount: perPageCount,
	}

	err := sp.Init(ctx)

	return &sp, err
}

// Single resolves the collection and applies the results to the given slice pointer.
func (sql SqlTable) Single(ctx context.Context, single interface{}) (err error) {
	selectQ := "select " + strings.Join(sql.asColumns, ", ") + " from " + sql.Table

	if len(sql.filterStatements) == 0 {
		err = sql.database.Select(single, selectQ+" limit 1")
	} else {
		err = sql.database.Select(single, selectQ+" where "+strings.Join(sql.filterStatements, " and ")+" limit 1", sql.filterValues...)
	}
	return err
}

// Delete resolves the collection using the filters and deletes the filtered elements.
func (sql SqlTable) Delete(ctx context.Context) (err error) {
	deleteQ := "delete from " + sql.Table

	if len(sql.filterStatements) == 0 {
		_, err = sql.database.Exec(deleteQ)
	} else {
		_, err = sql.database.Exec(deleteQ+" where "+strings.Join(sql.filterStatements, " and "), sql.filterValues...)
	}

	return err
}

// Update resolves the collection using the filters and updates the filtered elements using the operations
func (sql SqlTable) Update(ctx context.Context) (err error) {
	updateQ := "update " + sql.Table + " set "

	updateQ = updateQ + " " + strings.Join(sql.updateStatements, ", ")

	vals := make([]interface{}, 0)
	vals = append(vals, sql.updateValues...)
	vals = append(vals, sql.filterValues...)

	if len(sql.filterStatements) == 0 {
		_, err = sql.database.Exec(updateQ, vals...)
	} else {
		_, err = sql.database.Exec(updateQ+" where "+strings.Join(sql.filterStatements, " and "), vals...)
	}

	return err
}

// Insert inserts the given object
func (sql SqlTable) Insert(ctx context.Context, i interface{}) (err error) {
	insertQ := "insert into " + sql.Table + " (" + strings.Join(sql.insColumns, ",") + ") values "
	insertQ = insertQ + " (" + strings.Join(sql.namedColumns, ",") + ")"

	query, args, err := sqlx.Named(insertQ, i)

	arg2 := make([]interface{}, 0)

	for _, a := range args {
		if s, ok := a.(string); ok && s == "" {
			arg2 = append(arg2, nil)
		} else {
			arg2 = append(arg2, a)
		}
	}

	if err != nil {
		return err
	}

	_, err = sql.database.Exec(query, arg2...)

	return err
}

// ExecRaw executes a raw statement against the collection
func (sql SqlTable) ExecRaw(ctx context.Context, query string) error {
	_, err := sql.database.Exec(query)
	return err
}

// end Collection implementation
