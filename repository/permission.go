package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/google/uuid"
)

type Permissions struct {
	baseRepostiory
}

// Must be initiated after [User] and [Project]
func NewPermissions(
	ctx context.Context,
	db *sql.DB,
	log *slog.Logger,
	assert tinyssert.Assertions,
) (*Permissions, error) {
	b := newBaseRepostiory(ctx, db, log, assert)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS project_permissions (
		project_id        TEXT    NOT NULL,
		user_id	          TEXT    NOT NULL,
		permissions_value INTEGER NOT NULL DEFAULT '0',
		_permissions_text TEXT    NOT NULL DEFAULT '', -- For display purposes only, may not always be up-to-date
		created_at        TEXT    NOT NULL,
		updated_at        TEXT    NOT NULL,

		PRIMARY KEY(project_id, user_id)
		FOREIGN KEY(project_id)
			REFERENCES projects (id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT,
		FOREIGN KEY(user_id)
			REFERENCES users (id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT
	)	
	`)

	_, err = tx.ExecContext(ctx, q)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Join(errors.New("unable to create project tables"), err)
	}

	return &Permissions{baseRepostiory: b}, nil
}

func (repo Permissions) Create(project, user uuid.UUID, permissions model.Permissions) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.ctx)

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return errors.Join(ErrDatabaseConn, err)
	}

	q := `
	INSERT INTO project_permissions (project_id, user_id, permissions_value, _permissions_text, created_at, updated_at)
	VALUES (:project_id, :user_id, :permissions_value, :permissions_text, :created_at, :updated_at)
	`

	now := time.Now()

	log := repo.log.With(slog.String("project_id", project.String()),
		slog.String("user_id", user.String()),
		slog.String("permissions", fmt.Sprintf("%d", permissions)),
		slog.String("permissions_text", permissions.String()),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Inserting new project permissions")

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("project_id", project),
		sql.Named("user_id", user),
		sql.Named("permissions_value", permissions),
		sql.Named("permissions_text", permissions.String()),
		sql.Named("created_at", now.Format(dateFormat)),
		sql.Named("updated_at", now.Format(dateFormat)),
	)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to insert project permissions", slog.String("error", err.Error()))
		return errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return errors.Join(ErrCommitQuery, err)
	}

	return nil
}

func (repo Permissions) GetByID(project uuid.UUID, user uuid.UUID) (model.Permissions, error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.log)

	q := `
	SELECT permissions_value FROM project_permissions 
	WHERE project_id = :project_id
	AND   user_id    = :user_id
	`

	log := repo.log.With(slog.String("projcet_id", project.String()),
		slog.String("user_id", user.String()),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Getting by ID")

	row := repo.db.QueryRowContext(repo.ctx, q,
		sql.Named("project_id", user),
		sql.Named("user_id", user))

	var p model.Permissions
	if err := row.Scan(&p); err != nil {
		log.ErrorContext(repo.ctx, "Failed to get permissions by ID", slog.String("error", err.Error()))
		return model.Permissions(0), errors.Join(ErrExecuteQuery, err)
	}

	return p, nil
}

// GetByUserID returns a project_id-to-permissions map containing all projects and permissions that said userID
// has relation to.
func (repo Permissions) GetByUserID(user uuid.UUID) (permissions map[uuid.UUID]model.Permissions, err error) {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.log)

	// Begin tx so we don't read rows as they are being updated
	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return nil, errors.Join(ErrDatabaseConn, err)
	}

	q := `
	SELECT project_id, permissions_value FROM project_permissions 
	WHERE user_id = :user_id
	`

	log := repo.log.With(slog.String("user_id", user.String()),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Getting by user ID")

	rows, err := tx.QueryContext(repo.ctx, q, sql.Named("user_id", user))
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to get permissions by user ID", slog.String("error", err.Error()))
		return nil, errors.Join(ErrExecuteQuery, err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			err = errors.Join(ErrCloseConn, err)
		}
	}()

	ps := map[uuid.UUID]model.Permissions{}

	for rows.Next() {
		var project uuid.UUID
		var permissions model.Permissions

		err := rows.Scan(&project, &permissions)
		if err != nil {
			log.ErrorContext(repo.ctx, "Failed to scan permissions of user id", slog.String("error", err.Error()))
			return nil, errors.Join(ErrInvalidOutput, err)
		}

		ps[project] = permissions
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return nil, errors.Join(ErrCommitQuery, err)
	}

	return ps, nil
}

func (repo Permissions) Update(project, user uuid.UUID, permissions model.Permissions) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.log)

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return errors.Join(ErrDatabaseConn, err)
	}

	q := `
	UPDATE project_permissions
	SET permissions_value = :permissions_value
	    _permissions_text = :permissions_text
	    updated_at        = :updated_at
	WHERE project_uuid = :project_uuid
	  AND user_uuid    = :user_uuid
	`

	log := repo.log.With(slog.String("project_id", project.String()),
		slog.String("user_id", user.String()),
		slog.String("permissions", fmt.Sprintf("%d", permissions)),
		slog.String("permissions_text", permissions.String()),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Updating project permissions")

	now := time.Now()

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("permissions_value", permissions),
		sql.Named("permissions_text", permissions.String()),
		sql.Named("updated_at", now.Format(dateFormat)),
		sql.Named("project_id", project),
		sql.Named("user_id", user),
	)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to update project permissions", slog.String("error", err.Error()))
		return errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return errors.Join(ErrCommitQuery, err)
	}

	return nil
}

func (repo Permissions) Delete(project, user uuid.UUID) error {
	repo.assert.NotNil(repo.db)
	repo.assert.NotNil(repo.ctx)
	repo.assert.NotNil(repo.ctx)

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	if err != nil {
		return err
	}

	q := `
	DELETE FROM project_permissions 
	WHERE project_id = :project_id
	  AND user_id    = :user_id
	`

	log := repo.log.With(slog.String("project_id", project.String()),
		slog.String("user_id", user.String()),
		slog.String("query", q))
	log.DebugContext(repo.ctx, "Deleting project permissions")

	_, err = tx.ExecContext(repo.ctx, q,
		sql.Named("project_id", project),
		sql.Named("user_id", user),
	)
	if err != nil {
		log.ErrorContext(repo.ctx, "Failed to delete project permissions", slog.String("error", err.Error()))
		return errors.Join(ErrExecuteQuery, err)
	}

	if err := tx.Commit(); err != nil {
		log.ErrorContext(repo.ctx, "Failed to commit transaction", slog.String("error", err.Error()))
		return errors.Join(ErrCommitQuery, err)
	}

	return nil
}
