package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type Project struct {
	ID    string
	Title string
}

func (db *Database) setupProjects(tx *sql.Tx) error {
	db.assert.NotNil(tx)
	db.assert.NotNil(db.ctx)

	q := `CREATE TABLE IF NOT EXISTS projects (
	id    TEXT PRIMARY KEY NOT NULL,
	title TEXT NOT NULL
) STRICT`
	_, err := tx.ExecContext(db.ctx, q)
	if err != nil {
		return errors.Join(errors.New(`unable to execute create query to table "projects"`), err)
	}
	return nil
}

func (db *Database) CreateProject(id string, title string) (Project, error) {
	db.assert.NotNil(db.sql)
	db.assert.NotNil(db.ctx)
	db.assert.NotNil(db.log)
	db.assert.NotZero(id)
	db.assert.NotZero(title)

	q := fmt.Sprintf(`INSERT INTO projects (id, title) VALUES ('%s', '%s')`, id, title)

	db.log.Debug("Inserting into Projects", slog.String("query", q))

	tx, err := db.sql.BeginTx(db.ctx, nil)
	if err != nil {
		return Project{}, err
	}

	_, err = tx.ExecContext(db.ctx, q)
	if err != nil {
		return Project{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Project{}, err
	}

	return Project{ID: id, Title: title}, nil
}

func (db *Database) GetProject(id string) (Project, error) {
	db.assert.NotNil(db.sql)
	db.assert.NotNil(db.ctx)
	db.assert.NotNil(db.log)

	q := fmt.Sprintf(`SELECT id, title FROM projects WHERE id = '%s'`, id)

	db.log.Debug("Getting Project", slog.String("query", q))

	tx, err := db.sql.BeginTx(db.ctx, nil)
	if err != nil {
		return Project{}, err
	}

	row := tx.QueryRowContext(db.ctx, q)

	p := Project{}
	err = row.Scan(&p.ID, &p.Title)
	if err != nil {
		return p, err
	}

	err = tx.Commit()
	if err != nil {
		return p, err
	}

	return p, nil
}

func (db *Database) ListProjects() ([]Project, error) {
	db.assert.NotNil(db.sql)
	db.assert.NotNil(db.ctx)
	db.assert.NotNil(db.log)

	q := `SELECT id, title FROM projects`

	db.log.Debug("Listing Projects", slog.String("query", q))

	tx, err := db.sql.BeginTx(db.ctx, nil)
	if err != nil {
		return []Project{}, err
	}

	rows, err := tx.QueryContext(db.ctx, q)
	if err != nil {
		db.assert.Nil(tx.Rollback())
		return []Project{}, err
	}

	ps := []Project{}
	for rows.Next() {
		p := Project{}

		err := rows.Scan(&p.ID, &p.Title)
		if err != nil {
			db.assert.Nil(tx.Rollback())
			return ps, err
		}

		ps = append(ps, p)
	}

	err = tx.Commit()
	if err != nil {
		return ps, err
	}

	return ps, nil
}
