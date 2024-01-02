package database

import "context"

type IDatabase interface {
	ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) error
}

type ITransaction interface {
	ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) error
	Commit() error
	RollBack() error
}
