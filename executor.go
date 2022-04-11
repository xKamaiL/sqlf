package sqlf

import (
	"context"
	"database/sql"
	"reflect"

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

// StructElemToPtr get elem in struct to dest ptr
func StructElemToPtr(data any) []any {

	dest := make([]interface{}, 0)

	typ := reflect.TypeOf(data).Elem()
	val := reflect.ValueOf(data).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		t := typ.Field(i)
		if field.Kind() == reflect.Struct && t.Anonymous {
			StructElemToPtr(field.Addr().Interface())
		} else {
			dbFieldName := t.Tag.Get("db")
			if dbFieldName != "" {
				dest = append(dest, field.Addr().Interface())
			}
		}
	}
	return dest
}
func Iter[T any](ctx context.Context, q *Stmt, value any) error {
	return pgctx.Iter(ctx, func(scan pgsql.Scanner) error {
		var t T
		//q.To
		if err := scan(t); err != nil {
			return err
		}
		return nil
	}, q.String(), q.args...)
}
func (q *Stmt) Iter(ctx context.Context, f IteratorFunc) error {
	var err error
	err = pgctx.Iter(ctx, func(scan pgsql.Scanner) error {

		//var item T

		if err := scan(q.dest...); err != nil {
			return err
		}
		// TODO:

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
