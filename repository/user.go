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
)

type UserRepository struct {
	db *sql.DB

	ctx    context.Context
	log    *slog.Logger
	assert tinyssert.Assertions
}

func NewUserRepository(
	db *sql.DB,
	ctx context.Context,
	logger *slog.Logger,
	assert tinyssert.Assertions,
) (*UserRepository, error) {
	assert.NotNil(db)
	assert.NotNil(ctx)
	assert.NotNil(logger)

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		username      TEXT NOT NULL PRIMARY KEY,
		password_hash TEXT NOT NULL,
		created_at    TEXT NOT NULL,
		updated_at    TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return &UserRepository{
		db:     db,
		ctx:    ctx,
		log:    logger,
		assert: assert,
	}, nil
}

func (r *UserRepository) Create(u model.User) (model.User, error) {
	r.assert.NotNil(r.db)
	r.assert.NotNil(r.log)
	r.assert.NotNil(r.ctx)

	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return model.User{}, err
	}

	q := `
	INSERT INTO users (username, password_hash, created_at, updated_at)
	  VALUES (:username, :password_hash, :created_at, :updated_at)
	`

	log := r.log.With(slog.String("username", u.Username), slog.String("query", q))
	log.DebugContext(r.ctx, "Inserting new user")

	t := time.Now()

	passwd := base64.URLEncoding.EncodeToString(u.Password)

	_, err = tx.ExecContext(r.ctx, q,
		sql.Named("username", u.Username),
		sql.Named("password_hash", passwd),
		sql.Named("created_at", t.Format(dateFormat)),
		sql.Named("updated_at", t.Format(dateFormat)))
	if err != nil {
		log.ErrorContext(r.ctx, "Failed to create user", slog.String("error", err.Error()))
		return model.User{}, err
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(r.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return model.User{}, err
	}

	return u, nil
}

func (r *UserRepository) GetByUsername(username string) (model.User, error) {
	r.assert.NotNil(r.db)
	r.assert.NotNil(r.log)
	r.assert.NotNil(r.ctx)

	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return model.User{}, err
	}

	q := `
	SELECT username, password_hash, created_at, updated_at FROM users
	  WHERE username = :username
	`

	log := r.log.With(slog.String("username", username), slog.String("query", q))
	log.DebugContext(r.ctx, "Querying user")

	row := tx.QueryRowContext(r.ctx, q, sql.Named("username", username))

	var password_hash, dateCreated, dateUpdated string
	if err = row.Scan(&username, &password_hash, &dateCreated, &dateUpdated); err != nil {
		return model.User{}, err
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(r.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return model.User{}, err
	}

	passwd, err := base64.URLEncoding.DecodeString(password_hash)
	if err != nil {
		return model.User{}, err
	}

	c, err := time.Parse(dateFormat, dateCreated)
	if err != nil {
		return model.User{}, errors.Join(ErrInvalidData, err)
	}

	u, err := time.Parse(dateFormat, dateUpdated)
	if err != nil {
		return model.User{}, errors.Join(ErrInvalidData, err)
	}

	return model.User{
		Username:    username,
		Password:    passwd,
		DateCreated: c,
		DateUpdated: u,
	}, nil
}

func (r *UserRepository) Delete(u model.User) error {
	r.assert.NotNil(r.db)
	r.assert.NotNil(r.log)
	r.assert.NotNil(r.ctx)

	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}

	q := `
	DELETE FROM users WHERE username = :username
	`

	log := r.log.With(slog.String("username", u.Username), slog.String("query", q))
	log.DebugContext(r.ctx, "Deleting user")

	_, err = tx.ExecContext(r.ctx, q, sql.Named("username", u.Username))
	if err != nil {
		log.ErrorContext(r.ctx, "Failed to delete user", slog.String("error", err.Error()))
		return err
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(r.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return err
	}

	return nil
}
