package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/google/uuid"
)

type Project struct {
	baseRepostiory
}

func NewProject(ctx context.Context, db *sql.DB, log *slog.Logger, assert tinyssert.Assertions) (*Project, error) {
	b := newBaseRepostiory(ctx, db, log, assert)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS projects (
		id		  TEXT NOT NULL PRIMARY KEY,
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

	return &Project{baseRepostiory: b}, nil
}

func (repo Project) Create(p model.Project) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.ctx)

	if err := p.Validate(); err != nil {
		return errors.Join(ErrInvalidInput, err)
	}

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return errors.Join(ErrDatabaseConn, err)
	}

	q := `
	INSERT INTO users (id, title, created_at, updated_at)
	  VALUES (:id, :title, :created_at, :updated_at)
	`

	log := repo.log.With(slog.String("id", p.ID.String()), slog.String("query", q))
	log.DebugContext(repo.ctx, "Inserting new project")

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("id", p.ID),
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

func (repo Project) Update(p model.Project) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.ctx)

	if err := p.Validate(); err != nil {
		return errors.Join(ErrInvalidInput, err)
	}

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return errors.Join(ErrDatabaseConn, err)
	}

	q := `
	UPDATE projects 
	SET title = :title
	    updated_at = :updated_at
	WHERE id = :id
	`

	log := repo.log.With(slog.String("id", p.ID.String()), slog.String("query", q))
	log.DebugContext(repo.ctx, "Updating project")

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("title", p.Title),
		sql.Named("updated_at", p.DateUpdated.Format(dateFormat)),
		sql.Named("id", p.ID),
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

func (repo Project) DeleteByID(id uuid.UUID) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.ctx)

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return err
	}

	q := `
	DELETE FROM projects WHERE id = :id
	`

	log := repo.log.With(slog.String("id", id.String()), slog.String("query", q))
	log.DebugContext(repo.ctx, "Deleting project")

	_, err = tx.ExecContext(repo.ctx, q, sql.Named("id", id))
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
