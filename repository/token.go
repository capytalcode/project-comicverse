package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/google/uuid"
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

func (repo Token) Get(tokenID, userID uuid.UUID) (model.Token, error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.log)

	q := `
	SELECT (id, user_id, created_at, expired_at) FROM tokens
	WHERE id      = :id
	  AND user_id = :user_id
	`

	log := repo.log.With(slog.String("id", tokenID.String()),
		slog.String("user_id", userID.String()),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Getting token")

	row := repo.db.QueryRowContext(repo.ctx, q,
		sql.Named("id", tokenID),
		sql.Named("user_id", userID),
	)

	token, err := repo.scan(row)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to scan token", slog.String("error", err.Error()))
		return model.Token{}, err
	}

	return token, nil
}

func (repo Token) GetByUserID(userID uuid.UUID) ([]model.Token, error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.log)

	q := `
	SELECT (id, user_id, created_at, expired_at) FROM tokens
	WHERE user_id = :user_id
	`

	log := repo.log.With(
		slog.String("user_id", userID.String()),
		slog.String("query", q),
	)
	log.DebugContext(repo.ctx, "Getting users tokens")

	rows, err := repo.db.QueryContext(repo.ctx, q,
		sql.Named("user_id", userID),
	)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to get user tokens", slog.String("error", err.Error()))
		return []model.Token{}, errors.Join(ErrExecuteQuery, err)
	}

	tokens := []model.Token{}
	for rows.Next() {
		t, err := repo.scan(rows)
		if err != nil {
			log.ErrorContext(repo.ctx, "Failed to scan token", slog.String("error", err.Error()))
			return []model.Token{}, err
		}

		tokens = append(tokens, t)
	}

	if err := rows.Err(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to scan token rows", slog.String("error", err.Error()))
		return []model.Token{}, errors.Join(ErrExecuteQuery, err)
	}

	return tokens, nil
}

func (repo Token) scan(row scan) (model.Token, error) {
	repo.assert.NotNil(repo.ctx)

	var token model.Token
	var createdStr, expiresStr string

	err := row.Scan(&token.ID, &token.UserID, &createdStr, &expiresStr)
	if err != nil {
		return model.Token{}, errors.Join(ErrExecuteQuery, err)
	}

	dateCreated, err := time.Parse(dateFormat, createdStr)
	if err != nil {
		return model.Token{}, errors.Join(ErrInvalidOutput, err)
	}

	dateExpires, err := time.Parse(dateFormat, createdStr)
	if err != nil {
		return model.Token{}, errors.Join(ErrInvalidOutput, err)
	}

	token.DateCreated = dateCreated
	token.DateExpires = dateExpires

	if err := token.Validate(); err != nil {
		return model.Token{}, errors.Join(ErrInvalidOutput, err)
	}

	return token, nil
}

