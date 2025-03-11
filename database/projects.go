package database

import (
	"errors"
	"fmt"
)

type Project struct {
	ID    string
	Title string
}

var _ Table = (*Project)(nil)

func (p *Project) setup() string {
	return `
CREATE TABLE IF NOT EXISTS projects (
	id    TEXT PRIMARY KEY NOT NULL,
	title TEXT NOT NULL
) STRICT`
}

func (p *Project) insert() (string, error) {
	if p.ID == "" {
		return "", errors.New("ID field shouldn't be a empty string")
	}
	if p.Title == "" {
		return "", errors.New("Title field shouldn't be a empty string")
	}
	return fmt.Sprintf(`
INSERT OR FAIL INTO projects (
	id,
	title
) VALUES (
	'%s',
	'%s'
)`, p.ID, p.Title), nil
}
