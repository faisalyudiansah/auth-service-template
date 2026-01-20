package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/logger"
)

type Row struct {
	row       *sql.Row
	ctx       context.Context
	query     string
	args      []any
	startTime time.Time
	debug     bool
	slowLimit time.Duration
}

func (r *Row) Scan(dest ...any) error {
	err := r.row.Scan(dest...)
	logQuery(
		r.ctx,
		r.debug,
		r.slowLimit,
		r.query,
		r.args,
		err,
		time.Since(r.startTime),
	)
	return err
}

type DB struct {
	*sql.DB
	Debug     bool
	SlowLimit time.Duration
}

type Tx struct {
	*sql.Tx
	debug     bool
	slowLimit time.Duration
}

type Executor interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := db.DB.QueryContext(ctx, query, args...)
	logQuery(ctx, db.Debug, db.SlowLimit, query, args, err, time.Since(start))
	return rows, err
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *Row {
	return &Row{
		row:       db.DB.QueryRowContext(ctx, query, args...),
		ctx:       ctx,
		query:     query,
		args:      args,
		startTime: time.Now(),
		debug:     db.Debug,
		slowLimit: db.SlowLimit,
	}
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := db.DB.ExecContext(ctx, query, args...)
	logQuery(ctx, db.Debug, db.SlowLimit, query, args, err, time.Since(start))
	return res, err
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	logTransaction(ctx, "BEGIN", err)
	return &Tx{
		Tx:        tx,
		debug:     db.Debug,
		slowLimit: db.SlowLimit,
	}, nil
}

func (tx *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	logQuery(ctx, tx.debug, tx.slowLimit, query, args, err, time.Since(start))
	return rows, err
}

func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...any) *Row {
	return &Row{
		row:       tx.Tx.QueryRowContext(ctx, query, args...),
		ctx:       ctx,
		query:     query,
		args:      args,
		startTime: time.Now(),
		debug:     tx.debug,
		slowLimit: tx.slowLimit,
	}
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := tx.Tx.ExecContext(ctx, query, args...)
	logQuery(ctx, tx.debug, tx.slowLimit, query, args, err, time.Since(start))
	return res, err
}

func (tx *Tx) Commit(ctx context.Context) error {
	err := tx.Tx.Commit()
	logTransaction(ctx, "COMMIT", err)
	return err
}

func (tx *Tx) Rollback(ctx context.Context) error {
	err := tx.Tx.Rollback()
	logTransaction(ctx, "ROLLBACK", err)
	return err
}

func logQuery(ctx context.Context, debug bool, slowLimit time.Duration, query string, args []any, err error, duration time.Duration) {
	if !debug && duration < slowLimit && err == nil {
		return
	}

	level := "INFO"
	bg := "47" // default putih

	if duration >= slowLimit {
		level = "SLOW"
		bg = "43"
	}
	if err != nil {
		level = "ERROR"
		bg = "41"
	}

	renderedSQL := interpolate(query, args)

	colorEnd := "\033[0m" // reset
	blue := "\033[1;34m"

	msg := fmt.Sprintf(
		"%s %s%s%s\nduration=%s err=%v",
		bgTag(bg, "[SQL]["+level+"]"),
		blue,
		renderedSQL,
		colorEnd,
		duration,
		err,
	)

	logger.PushQueryAsyncLog(ctx, msg)
}

func logTransaction(ctx context.Context, action string, err error) {
	bg := "46" // default putih
	status := "SUCCESS"

	if err != nil {
		bg = "41"
		status = "FAILED"
	} else {
		switch action {
		case "COMMIT":
			bg = "42"
		case "ROLLBACK":
			bg = "43"
		case "BEGIN":
			bg = "46"
		}
	}

	msg := fmt.Sprintf(
		"%s err=%v\n",
		bgTag(bg, "[TRANSACTION-"+action+"]["+status+"]"),
		err,
	)
	logger.PushQueryAsyncLog(ctx, msg)
}

func bgTag(bgColor string, text string) string {
	fg := "37"
	if bgColor == "47" {
		fg = "30"
	}
	return fmt.Sprintf("\033[%s;%sm%s\033[0m", fg, bgColor, text)
}

func interpolate(query string, args []any) string {
	q := query

	for i, arg := range args {
		ph := fmt.Sprintf("$%d", i+1)

		var val string
		switch v := arg.(type) {
		case string:
			val = "'" + strings.ReplaceAll(v, "'", "''") + "'"
		case time.Time:
			val = "'" + v.Format(time.RFC3339) + "'"
		case nil:
			val = "NULL"
		default:
			val = fmt.Sprintf("%v", v)
		}

		q = strings.Replace(q, ph, val, 1)
	}

	return compact(q)
}

func compact(q string) string {
	return strings.Join(strings.Fields(q), " ")
}
