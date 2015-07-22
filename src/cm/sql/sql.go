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

type SqlTable struct {
	Table string

	database         *sqlx.DB
	columns          []string
	filterStatements []string
	values           []interface{}
}

func Column(name string, ctype string) cm.ValueColumn {
	return SqlValueColumn{name, ctype}
}

type SqlValueColumn struct {
	name  string
	ctype string
}

type SqlEqPredicate struct {
	Column cm.ValueColumn
	Value  interface{}
}

type SqlNotEqPredicate struct {
	Column cm.ValueColumn
	Value  interface{}
}

type SqlLikePredicate struct {
	Column        cm.ValueColumn
	Value         interface{}
	CaseSensitive bool
}

func (s SqlValueColumn) Name() string {
	return s.name
}

func (s SqlValueColumn) Type() string {
	return s.ctype
}

func (s SqlValueColumn) Eq(i interface{}) cm.Predicate {
	return &SqlEqPredicate{s, i}
}

func (s SqlValueColumn) NotEq(i interface{}) cm.Predicate {
	return &SqlNotEqPredicate{s, i}
}

func (s SqlValueColumn) Like(caseSensitive bool, i interface{}) cm.Predicate {
	return &SqlLikePredicate{s, i, caseSensitive}
}

func (pred *SqlEqPredicate) Apply(c cm.Collection) error {
	col := c.(*SqlTable)
	col.filterStatements = append(col.filterStatements,
		fmt.Sprintf("%s = ?", pred.Column.Name()))

	col.values = append(col.values, pred.Value)

	return nil
}

func (pred *SqlNotEqPredicate) Apply(c cm.Collection) error {
	col := c.(*SqlTable)
	col.filterStatements = append(col.filterStatements,
		fmt.Sprintf("%s != ?", pred.Column.Name()))

	col.values = append(col.values, pred.Value)

	return nil
}

func (pred *SqlLikePredicate) Apply(c cm.Collection) error {
	col := c.(*SqlTable)

	like := "like"

	if pred.CaseSensitive {
		like = "ilike"
	}

	col.filterStatements = append(col.filterStatements,
		fmt.Sprintf("%s %s ?", pred.Column.Name(), like))

	col.values = append(col.values, pred.Value)

	return nil
}

func New(db *sqlx.DB, table string) *SqlTable {
	return &SqlTable{
		Table:            table,
		database:         db,
		columns:          make([]string, 0),
		filterStatements: make([]string, 0),
		values:           make([]interface{}, 0),
	}
}

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

func (sql SqlTable) Filter(pred cm.Predicate) cm.Collection {
	pred.Apply(&sql)
	return &sql
}

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

func (sql SqlTable) Single(ctx context.Context, single interface{}) (err error) {
	selectQ := "select " + strings.Join(sql.columns, ", ") + " from " + sql.Table

	if len(sql.filterStatements) == 0 {
		err = sql.database.Select(single, selectQ+" limit 1")
	} else {
		err = sql.database.Select(single, selectQ+" where "+strings.Join(sql.filterStatements, " and ")+" limit 1", sql.values...)
	}
	return err
}

func (sql SqlTable) Delete(ctx context.Context) (err error) {
	deleteQ := "delete from " + sql.Table

	if len(sql.filterStatements) == 0 {
		_, err = sql.database.Exec(deleteQ)
	} else {
		_, err = sql.database.Exec(deleteQ+" where "+strings.Join(sql.filterStatements, " and "), sql.values...)
	}

	return err
}

func (sql SqlTable) ExecRaw(ctx context.Context, query string) error {
	_, err := sql.database.Exec(query)
	return err
}
