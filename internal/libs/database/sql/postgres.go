package sql

import (
	"context"
	"database/sql"
)

type ISqlDatabase interface {
	ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) (string, error)
}

type ISqlTransaction interface {
	ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) (string, error)
	Commit() error
	RollBack() error
}

type database struct {
	conn *sql.DB
}

type transaction struct {
	txn *sql.Tx
}

func NewDatabase(connectionString string) (ISqlDatabase, error) {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &database{
		conn: conn,
	}, nil
}

func (d *database) Stop() error {
	err := d.conn.Close()
	return err
}

func (d *database) ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) (string, error) {
	rows, err := d.conn.QueryContext(ctx, queryString)
	if err != nil {
		return "", err
	}
	if len(opt) == 0 {
		return "", nil
	}
	err = rows.Scan(opt[0])
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (d *database) Transaction(ctx context.Context) (ISqlTransaction, error) {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction{
		txn: tx,
	}, nil
}

func (t *transaction) ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) (string, error) {
	rows, err := t.txn.QueryContext(ctx, queryString)
	if err != nil {
		return "", err
	}
	if len(opt) == 0 {
		return "", nil
	}
	err = rows.Scan(opt[0])
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (t *transaction) Commit() error {
	if err := t.txn.Commit(); err != nil {
		return err
	}
	return nil
}

func (t *transaction) RollBack() error {
	if err := t.txn.Rollback(); err != nil {
		return err
	}
	return nil
}
