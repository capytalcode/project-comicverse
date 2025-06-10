package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"forge.capytal.company/loreddev/x/tinyssert"
)

type Token struct {
	baseRepostiory
}

// Must be initiated after [User]
func NewToken(ctx context.Context, db *sql.DB, log *slog.Logger, assert tinyssert.Assertions) (*Token, error) {
	b := newBaseRepostiory(ctx, db, log, assert)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS tokens (
		id	       TEXT NOT NULL,
		user_id    TEXT NOT NULL,
		created_at TEXT NOT NULL,
		expires_at TEXT NOT NULL,

		PRIMARY KEY(uuid, user_id),
		FOREIGN KEY(user_id)
			REFERENCES users (id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT
	)`)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Join(errors.New("unable to create project tables"), err)
	}

	return &Token{baseRepostiory: b}, nil
}

