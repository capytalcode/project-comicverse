package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"forge.capytal.company/capytalcode/project-comicverse/model"
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

func (repo ProjectRepository) Create(p model.Project) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.ctx)

	if err := p.Validate(); err != nil {
		return errors.Join(ErrInvalidData, err)
	}

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return errors.Join(ErrDatabaseConn, err)
	}

	q := `
	INSERT INTO users (uuid, title, created_at, updated_at)
	  VALUES (:uuid, :title, :created_at, :updated_at)
	`

	log := repo.log.With(slog.String("uuid", p.UUID.String()), slog.String("query", q))
	log.DebugContext(repo.ctx, "Inserting new project")

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("uuid", p.UUID),
		sql.Named("title", p.Title),
		sql.Named("created_at", p.DateCreated.Format(dateFormat)),
		sql.Named("updated_at", p.DateUpdated.Format(dateFormat)),
	)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to insert project", slog.String("error", err.Error()))
		return errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return errors.Join(ErrCommitQuery, err)
	}

	return nil
}

func (repo ProjectRepository) Update(p model.Project) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.ctx)

	if err := p.Validate(); err != nil {
		return errors.Join(ErrInvalidData, err)
	}

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return errors.Join(ErrDatabaseConn, err)
	}

	q := `
	UPDATE projects 
	SET title = :title
	    updated_at = :updated_at
	WHERE uuid = :uuid
	`

	log := repo.log.With(slog.String("uuid", p.UUID.String()), slog.String("query", q))
	log.DebugContext(repo.ctx, "Updating project")

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("title", p.Title),
		sql.Named("updated_at", p.DateUpdated.Format(dateFormat)),
	)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to insert project", slog.String("error", err.Error()))
		return errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return errors.Join(ErrCommitQuery, err)
	}

	return nil
}

func (repo ProjectRepository) Delete(p model.Project) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.ctx)

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return err
	}

	q := `
	DELETE FROM projects WHERE uuid = :uuid
	`

	log := repo.log.With(slog.String("uuid", p.UUID.String()), slog.String("query", q))
	log.DebugContext(repo.ctx, "Deleting project")

	_, err = tx.ExecContext(repo.ctx, q, sql.Named("uuid", p.UUID))
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to delete project", slog.String("error", err.Error()))
		return errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return errors.Join(ErrCommitQuery, err)
	}

	return nil
}
