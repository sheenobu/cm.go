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

// Table defines the Collection that interacts with a sqlx.DB connection
type Table struct {
	Name string

	database      *sqlx.DB
	asColumnTypes []ValueColumn
	asColumns     []string
	insColumns    []string
	namedColumns  []string

	offset int

	filterStatements []string
	updateStatements []string
	updateValues     []interface{}
	filterValues     []interface{}
}

// New defines a new Table based on the database connection and table name.
func New(db *sqlx.DB, table string) *Table {
	return &Table{
		Name:     table,
		database: db,
	}
}

// begin Collection implementation

// Init initializes the sql table by performing reflection operations
// on the interface
func (sql *Table) Init(iface interface{}) error {
	var columns []string

	s := reflect.ValueOf(iface).Elem()
	st := reflect.TypeOf(iface).Elem()

	for i := 0; i < st.NumField(); i++ {
		f := s.Field(i)
		fx := st.Field(i)

		if f.Type().Name() == "ValueColumn" {

			vc := f.Interface().(ValueColumn)

			columns = append(columns, vc.Build())

			cs := camelcase.Split(fx.Name)

			for i := 0; i != len(cs); i++ {
				cs[i] = strings.ToLower(cs[i])
			}

			sql.asColumnTypes = append(sql.asColumnTypes, vc)
			sql.asColumns = append(sql.asColumns,
				vc.Name()+" as "+strings.ToLower(fx.Name))

			if !strings.Contains(vc.Type(), "AUTOINCREMENT") {
				sql.insColumns = append(sql.insColumns,
					vc.Name())

				sql.namedColumns = append(sql.namedColumns,
					":"+strings.ToLower(fx.Name))
			} else {
				sql.offset++
			}

		}
	}

	createTable := fmt.Sprintf(`
		create table %s
			(%s);
	`, sql.Name, strings.Join(columns, ", "))

	_, err := sql.database.Exec(createTable)

	return mapError(err)
}

// Filter filters the collection on the given predicate
func (sql Table) Filter(pred cm.Predicate) cm.Collection {
	pred.Apply(&sql)
	return &sql
}

// Edit updates the collection for the given operation
func (sql Table) Edit(op cm.Operation) cm.Collection {
	op.Apply(&sql)
	return &sql
}

// List resolves the collection and applies the results to the given slice pointer
func (sql Table) List(ctx context.Context, list interface{}) (err error) {

	selectQ := "select " + strings.Join(sql.asColumns, ", ") + " from " + sql.Name

	if len(sql.filterStatements) == 0 {
		err = sql.database.Select(list, selectQ)
	} else {
		q := selectQ + " where " + strings.Join(sql.filterStatements, " and ")
		err = sql.database.Select(list, q, sql.filterValues...)
	}

	return err
}

// Page returns a pagination object which is responsible for resolving the collection
func (sql Table) Page(ctx context.Context, perPageCount int) (cm.Paginator, error) {

	countQ := "select count(" + sql.asColumnTypes[0].name + ") from " + sql.Name
	selectQ := "select " + strings.Join(sql.asColumns, ", ") + " from " + sql.Name

	if len(sql.filterStatements) != 0 {
		selectQ = selectQ + " where " + strings.Join(sql.filterStatements, " and ")
		countQ = countQ + " where " + strings.Join(sql.filterStatements, " and ")
	}

	sp := Paginator{
		database:     sql.database,
		selectQuery:  selectQ,
		countQuery:   countQ,
		perPageCount: perPageCount,
	}

	err := sp.Init(ctx)

	return &sp, err
}

// Single resolves the collection and applies the results to the given slice pointer.
func (sql Table) Single(ctx context.Context, single interface{}) (err error) {
	selectQ := "select " + strings.Join(sql.asColumns, ", ") + " from " + sql.Name

	if len(sql.filterStatements) == 0 {
		err = sql.database.Select(single, selectQ+" limit 1")
	} else {
		err = sql.database.Select(single, selectQ+" where "+strings.Join(sql.filterStatements, " and ")+" limit 1", sql.filterValues...)
	}
	return err
}

// Delete resolves the collection using the filters and deletes the filtered elements.
func (sql Table) Delete(ctx context.Context) (err error) {
	deleteQ := "delete from " + sql.Name

	if len(sql.filterStatements) == 0 {
		_, err = sql.database.Exec(deleteQ)
	} else {
		_, err = sql.database.Exec(deleteQ+" where "+strings.Join(sql.filterStatements, " and "), sql.filterValues...)
	}

	return err
}

// Update resolves the collection using the filters and updates the filtered elements using the operations
func (sql Table) Update(ctx context.Context) (err error) {
	updateQ := "update " + sql.Name + " set "

	updateQ = updateQ + " " + strings.Join(sql.updateStatements, ", ")

	var vals []interface{}
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
func (sql Table) Insert(ctx context.Context, i interface{}) (err error) {
	insertQ := "insert into " + sql.Name + " (" + strings.Join(sql.insColumns, ",") + ") values "
	insertQ = insertQ + " (" + strings.Join(sql.namedColumns, ",") + ")"

	query, args, err := sqlx.Named(insertQ, i)

	var arg2 []interface{}

	for idx, a := range args {

		vc := sql.asColumnTypes[idx+sql.offset]
		if ins, ok := vc.fns["insert"]; ok && ins != nil {
			b := ins()
			if b != nil {
				arg2 = append(arg2, b)
				continue
			}
		}

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
func (sql Table) ExecRaw(ctx context.Context, query string) error {
	_, err := sql.database.Exec(query)
	return err
}

// end Collection implementation
