package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"forge.capytal.company/loreddev/x/tinyssert"
)

type ProjectRepository struct {
	baseRepostiory
}

func NewProjectRepository(ctx context.Context, db *sql.DB, log *slog.Logger, assert tinyssert.Assertions) (*ProjectRepository, error) {
	b := newBaseRepostiory(ctx, db, log, assert)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS projects (
		uuid		  TEXT NOT NULL PRIMARY KEY,
		title         TEXT NOT NULL,
		created_at    TEXT NOT NULL,
		updated_at    TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Join(errors.New("unable to create project tables"), err)
	}

	return &ProjectRepository{baseRepostiory: b}, nil
}

