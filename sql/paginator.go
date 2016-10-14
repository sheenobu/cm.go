package sql

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

// Paginator is the SQL implementation of the cm.Paginator
type Paginator struct {
	database     *sqlx.DB
	selectQuery  string
	countQuery   string
	perPageCount int
	parent       Table

	// mutable state
	pageCount   int
	currentPage int
	itemCount   int
}

// Init initializes the paginator by running the initial
// count query
func (p *Paginator) Init(ctx context.Context) (err error) {

	var count []int
	err = p.database.Select(&count, p.countQuery)

	if len(count) == 0 {
		return fmt.Errorf("Query failed: %s", err)
	}

	p.itemCount = count[0]
	p.pageCount = count[0] / p.perPageCount

	remaining := count[0] % p.perPageCount
	if remaining != 0 {
		p.pageCount++
	}

	p.currentPage = 0

	return err
}

// begin cm.Paginator implementation

// PageCount returns the number of pages of this paginator
func (p *Paginator) PageCount() int {
	return p.pageCount
}

// CurrentPage returns the current page of the pagintor
func (p *Paginator) CurrentPage() int {
	return p.currentPage
}

// ItemCount returns the totel item count of the paginator
func (p *Paginator) ItemCount() int {
	return p.itemCount
}

// PerPageCount returns the per-page count of the pagintor
func (p *Paginator) PerPageCount() int {
	return p.perPageCount
}

// Next increases the current page by one, unless the current page
// is the last page.
func (p *Paginator) Next() bool {
	if p.currentPage >= p.pageCount-1 {
		p.currentPage = p.pageCount - 1
		if p.currentPage < 0 {
			p.currentPage = 0
		}
		return false
	}

	p.currentPage++

	return true
}

// Prev decreases the current page by one, unless the current page
// is 0
func (p *Paginator) Prev() bool {
	if p.currentPage <= 0 {
		p.currentPage = 0
		return false
	}

	p.currentPage--
	return true
}

// Apply runs the paginating query against the database, inserting it into the
// given list
func (p *Paginator) Apply(ctx context.Context, list interface{}) (err error) {
	db := getDb(ctx, p.database)
	q := p.selectQuery + " limit " + strconv.Itoa(p.perPageCount) + " offset " + strconv.Itoa(p.currentPage*p.perPageCount)
	err = db.Select(list, q, p.parent.filterValues...)
	return err
}

// end cm.Paginator implementation
