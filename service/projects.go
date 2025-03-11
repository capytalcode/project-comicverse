package service

import (
	"encoding/xml"
	"errors"

	"forge.capytal.company/capytalcode/project-comicverse/database"
)

func (s *service) NewProject() error {
	s.assert.NotNil(s.db)

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	s.assert.NotZero(id.String(), "UUID should never be invalid")

	err = s.db.Insert(&database.Project{
		ID:    id.String(),
		Title: "New Project",
	})
	if err != nil {
		return errors.Join(errors.New("database returned error while inserting new project"), err)
	}


	return nil
}

func (s *service) ListProjects() {
}
