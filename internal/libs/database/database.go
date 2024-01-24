package database

import "context"

type IDatabase interface {
	ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) (string, error)
}

type ITransaction interface {
	ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) (string, error)
	Commit() error
	RollBack() error
}
