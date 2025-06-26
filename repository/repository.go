package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"forge.capytal.company/loreddev/x/tinyssert"
)

// TODO: Add rowback to all return errors, or use context to cancel operations

type baseRepostiory struct {
	db *sql.DB

	ctx    context.Context
	log    *slog.Logger
	assert tinyssert.Assertions
}

func newBaseRepostiory(ctx context.Context, db *sql.DB, log *slog.Logger, assert tinyssert.Assertions) baseRepostiory {
	assert.NotNil(db)
	assert.NotNil(ctx)
	assert.NotNil(log)

	return baseRepostiory{
		db:     db,
		ctx:    ctx,
		log:    log,
		assert: assert,
	}
}

var (
	// TODO: Change all ErrDatabaseConn to ErrCloseConn
	ErrDatabaseConn  = errors.New("repository: failed to begin transaction/connection with database")
	ErrCloseConn     = errors.New("repository: failed to close/commit connection")
	ErrExecuteQuery  = errors.New("repository: failed to execute query")
	ErrCommitQuery   = errors.New("repository: failed to commit transaction")
	ErrInvalidInput  = errors.New("repository: data sent to save is invalid")
	ErrInvalidOutput = errors.New("repository: data found is not valid")
	ErrNotFound      = sql.ErrNoRows
)

var dateFormat = time.RFC3339

type scan interface {
	Scan(dest ...any) error
}
