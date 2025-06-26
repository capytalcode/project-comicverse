package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

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
		id		   TEXT NOT NULL PRIMARY KEY,
		title      TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
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
	INSERT INTO projects (id, title, created_at, updated_at)
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

func (repo Project) GetByID(projectID uuid.UUID) (project model.Project, err error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.log)

	q := `
	SELECT id, title, created_at, updated_at FROM projects
		WHERE id = :id
	`

	log := repo.log.With(slog.String("query", q), slog.String("id", projectID.String()))
	log.DebugContext(repo.ctx, "Getting project by ID")

	row := repo.db.QueryRowContext(repo.ctx, q, sql.Named("id", projectID))

	var id uuid.UUID
	var title string
	var dateCreatedStr, dateUpdatedStr string

	err = row.Scan(&id, &title, &dateCreatedStr, &dateUpdatedStr)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to scan projects with IDs", slog.String("error", err.Error()))
		return model.Project{}, errors.Join(ErrInvalidOutput, err)
	}

	dateCreated, err := time.Parse(dateFormat, dateCreatedStr)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to scan projects with IDs", slog.String("error", err.Error()))
		return model.Project{}, errors.Join(ErrInvalidOutput, err)
	}

	dateUpdated, err := time.Parse(dateFormat, dateUpdatedStr)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to scan projects with IDs", slog.String("error", err.Error()))
		return model.Project{}, errors.Join(ErrInvalidOutput, err)
	}

	return model.Project{
		ID:          id,
		Title:       title,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}, nil
}

func (repo Project) GetByIDs(ids []uuid.UUID) (projects []model.Project, err error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.log)

	// Begin tx so we don't read rows as they are being updated
	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return nil, errors.Join(ErrDatabaseConn, err)
	}

	c := make([]string, len(ids))
	for i, id := range ids {
		c[i] = fmt.Sprintf("id = '%s'", id.String())
	}

	q := fmt.Sprintf(`
	SELECT id, title, created_at, updated_at FROM projects
	WHERE %s
	`, strings.Join(c, " OR "))

	log := repo.log.With(slog.String("query", q))
	log.DebugContext(repo.ctx, "Getting projects by IDs")

	rows, err := tx.QueryContext(repo.ctx, q)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to get projects by IDs", slog.String("error", err.Error()))
		return nil, errors.Join(ErrExecuteQuery, err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			err = errors.Join(ErrCloseConn, err)
		}
	}()

	ps := []model.Project{}

	for rows.Next() {
		var id uuid.UUID
		var title string
		var dateCreatedStr, dateUpdatedStr string

		err := rows.Scan(&id, &title, &dateCreatedStr, &dateUpdatedStr)
		if err != nil {
			log.ErrorContext(repo.ctx, "Failed to scan projects with IDs", slog.String("error", err.Error()))
			return nil, errors.Join(ErrInvalidOutput, err)
		}

		dateCreated, err := time.Parse(dateFormat, dateCreatedStr)
		if err != nil {
			log.ErrorContext(repo.ctx, "Failed to scan projects with IDs", slog.String("error", err.Error()))
			return nil, errors.Join(ErrInvalidOutput, err)
		}

		dateUpdated, err := time.Parse(dateFormat, dateUpdatedStr)
		if err != nil {
			log.ErrorContext(repo.ctx, "Failed to scan projects with IDs", slog.String("error", err.Error()))
			return nil, errors.Join(ErrInvalidOutput, err)
		}

		ps = append(ps, model.Project{
			ID:          id,
			Title:       title,
			DateCreated: dateCreated,
			DateUpdated: dateUpdated,
		})
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return nil, errors.Join(ErrCommitQuery, err)
	}

	return ps, nil
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
