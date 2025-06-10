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

func (repo Token) Create(token model.Token) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.log)

	if err := token.Validate(); err != nil {
		return errors.Join(ErrInvalidInput, err)
	}

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return errors.Join(ErrDatabaseConn, err)
	}

	q := `
	INSERT INTO tokens (id, user_id, created_at, expires_at)
	VALUES (:id, :user_id, :created_at, :expires_at)
	`

	log := repo.log.With(slog.String("id", token.ID.String()),
		slog.String("user_id", token.UserID.String()),
		slog.String("expires", token.DateExpires.Format(dateFormat)),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Inserting new user token")

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("id", token.ID),
		sql.Named("user_id", token.UserID),
		sql.Named("created_at", token.DateCreated.Format(dateFormat)),
		sql.Named("expired_at", token.DateExpires.Format(dateFormat)),
	)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to insert token", slog.String("error", err.Error()))
		return errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return errors.Join(ErrCommitQuery, err)
	}

	return nil
}

