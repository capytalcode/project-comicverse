package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/google/uuid"
)

type UserRepository struct {
	baseRepostiory
}

func NewUserRepository(
	ctx context.Context,
	db *sql.DB,
	logger *slog.Logger,
	assert tinyssert.Assertions,
) (*UserRepository, error) {
	assert.NotNil(ctx)
	assert.NotNil(db)
	assert.NotNil(logger)

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id            TEXT NOT NULL PRIMARY KEY,
		username      TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at    TEXT NOT NULL,
		updated_at    TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	b := newBaseRepostiory(ctx, db, logger, assert)

	return &UserRepository{
		baseRepostiory: b,
	}, nil
}

func (repo *UserRepository) Create(u model.User) (model.User, error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.log)
	repo.assert.NotNil(repo.ctx)

	if err := u.Validate(); err != nil {
		return model.User{}, err
	}

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return model.User{}, errors.Join(ErrDatabaseConn, err)
	}

	q := `
	INSERT INTO users (id, username, password_hash, created_at, updated_at)
	VALUES (:id, :username, :password_hash, :created_at, :updated_at)
	`

	log := repo.log.With(
		slog.String("id", u.ID.String()),
		slog.String("username", u.Username),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Inserting new user")

	t := time.Now()

	passwd := base64.URLEncoding.EncodeToString(u.Password)

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("id", u.ID),
		sql.Named("username", u.Username),
		sql.Named("password_hash", passwd),
		sql.Named("created_at", t.Format(dateFormat)),
		sql.Named("updated_at", t.Format(dateFormat)))
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to create user", slog.String("error", err.Error()))
		return model.User{}, errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return model.User{}, errors.Join(ErrCommitQuery, err)
	}

	return u, nil
}

func (repo *UserRepository) GetByID(id uuid.UUID) (model.User, error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.log)
	repo.assert.NotNil(repo.ctx)

	q := `
	SELECT id, username, password_hash, created_at, updated_at FROM users
	  WHERE id = :id
	`

	log := repo.log.With(
		slog.String("id", id.String()),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Querying user")

	row := repo.db.QueryRowContext(repo.ctx, q, sql.Named("username", id))

	user, err := repo.scan(row)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to query user", slog.String("error", err.Error()))
		return model.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) GetByUsername(username string) (model.User, error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.log)
	repo.assert.NotNil(repo.ctx)

	q := `
	SELECT id, username, password_hash, created_at, updated_at FROM users
	  WHERE username = :username
	`

	log := repo.log.With(
		slog.String("username", username),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Querying user")

	row := repo.db.QueryRowContext(repo.ctx, q, sql.Named("username", username))

	user, err := repo.scan(row)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to query user", slog.String("error", err.Error()))
		return model.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) scan(row scan) (model.User, error) {
	var user model.User
	var password_hashStr, createdStr, updatedStr string
	err := row.Scan(&user.ID, &user.Username, &password_hashStr, &createdStr, &updatedStr)
	if err != nil {
		return model.User{}, errors.Join(ErrExecuteQuery, err)
	}

	passwd, err := base64.URLEncoding.DecodeString(password_hashStr)
	if err != nil {
		return model.User{}, errors.Join(ErrInvalidOutput, err)
	}

	created, err := time.Parse(dateFormat, createdStr)
	if err != nil {
		return model.User{}, errors.Join(ErrInvalidOutput, err)
	}

	updated, err := time.Parse(dateFormat, updatedStr)
	if err != nil {
		return model.User{}, errors.Join(ErrInvalidOutput, err)
	}

	user.Password = passwd
	user.DateCreated = created
	user.DateUpdated = updated

	if err := user.Validate(); err != nil {
		return model.User{}, errors.Join(ErrInvalidOutput, err)
	}

	return user, nil
}

func (repo *UserRepository) Delete(u model.User) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.log)
	repo.assert.NotNil(repo.ctx)

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return err
	}

	q := `
	DELETE FROM users WHERE id = :id
	`

	log := repo.log.With(slog.String("id", u.ID.String()), slog.String("query", q))
	log.DebugContext(repo.ctx, "Deleting user")

	_, err = tx.ExecContext(repo.ctx, q, sql.Named("id", u.ID))
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to delete user", slog.String("error", err.Error()))
		return errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return errors.Join(ErrCommitQuery, err)
	}

	return nil
}
