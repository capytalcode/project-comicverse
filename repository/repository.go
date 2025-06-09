package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"forge.capytal.company/loreddev/x/tinyssert"
)

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
	ErrDatabaseConn = errors.New("failed to begin transaction/connection with database")
	ErrExecuteQuery = errors.New("failed to execute query")
	ErrCommitQuery  = errors.New("failed to commit transaction")
	ErrInvalidData  = errors.New("data sent to save is invalid")
	ErrNotFound     = sql.ErrNoRows
)

var dateFormat = time.RFC3339
