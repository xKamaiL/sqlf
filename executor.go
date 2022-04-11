package sqlf

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"
)

// Executor performs SQL queries.
// It's an interface accepted by Query, QueryRow and Exec methods.
// Both sql.DB, sql.Conn and sql.Tx can be passed as executor.
type Executor interface {
	pgctx.DB
}
type IteratorFunc func()

// QueryRow executes the statement via Executor methods
// and scans values to variables bound via To method calls.
// and call to pgctx.QueryRow
func (q *Stmt) QueryRow(ctx context.Context) error {
	return pgctx.QueryRow(ctx, q.String(), q.args...).Scan(q.dest...)
}

// Exec executes the statement.
func (q *Stmt) Exec(ctx context.Context) (sql.Result, error) {
	return pgctx.Exec(ctx, q.String(), q.args...)
}

func (q *Stmt) Iter(ctx context.Context, f IteratorFunc) error {
	var err error
	err = pgctx.Iter(ctx, func(scan pgsql.Scanner) error {
		if err := scan(q.dest...); err != nil {
			return err
		}
		f()
		return nil
	}, q.String(), q.args...)
	return err
}

// ExecAndClose executes the statement and releases all the objects
// and buffers allocated by statement builder back to a pool.
//
// Do not call any Stmt methods after this call.
func (q *Stmt) ExecAndClose(ctx context.Context) (sql.Result, error) {
	res, err := q.Exec(ctx)
	q.Close()
	return res, err
}
