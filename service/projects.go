package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"forge.capytal.company/capytalcode/project-comicverse/database"
	"forge.capytal.company/capytalcode/project-comicverse/internals/randstr"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const projectIDLength = 6

var ErrProjectNotExists = errors.New("project does not exists in database")

func (s *Service) CreateProject() (Project, error) {
	s.assert.NotNil(s.db)
	s.assert.NotNil(s.s3)
	s.assert.NotNil(s.ctx)
	s.assert.NotZero(s.bucket)

	s.log.Debug("Creating new project")

	id, err := randstr.NewHex(projectIDLength)
	if err != nil {
		return Project{}, errors.Join(errors.New("creating hexadecimal ID returned error"), err)
	}

	title := "New Project"

	s.assert.NotZero(id, "ID should never be empty")

	s.log.Debug("Creating project on database", slog.String("id", id))

	_, err = s.db.CreateProject(id, title)
	if err != nil {
		return Project{}, err
	}

	p := Project{
		ID:    id,
		Title: title,
		Pages: map[string]ProjectPage{},
	}

	c, err := json.Marshal(p)
	if err != nil {
		return Project{}, err
	}

	s.log.Debug("Creating project on storage", slog.String("id", id))

	f := fmt.Sprintf("%s.comic.json", id)
	_, err = s.s3.PutObject(s.ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &f,
		Body:   bytes.NewReader(c),
	})
	if err != nil {
		return Project{}, err
	}

	return p, nil
}

func (s *Service) GetProject(id string) (Project, error) {
	s.assert.NotNil(s.db)
	s.assert.NotNil(s.s3)
	s.assert.NotZero(s.bucket)
	s.assert.NotNil(s.ctx)
	s.assert.NotZero(id)

	res, err := s.db.GetProject(id)
	if errors.Is(err, database.ErrNoRows) {
		return Project{}, errors.Join(ErrProjectNotExists, err)
	}
	if err != nil {
		return Project{}, err
	}

	f := fmt.Sprintf("%s.comic.json", id)
	file, err := s.s3.GetObject(s.ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &f,
	})
	if err != nil {
		return Project{}, err
	}

	c, err := io.ReadAll(file.Body)
	if err != nil {
		return Project{}, err
	}

	var p Project
	err = json.Unmarshal(c, &p)

	s.assert.Equal(res.ID, p.ID, "The project ID should always be equal in the Database and Storage")
	s.assert.Equal(res.Title, p.Title)

	return p, err
}

func (s *Service) ListProjects() ([]Project, error) {
	s.assert.NotNil(s.db)

	ps, err := s.db.ListProjects()
	if err != nil {
		return []Project{}, err
	}

	p := make([]Project, len(ps))
	for i, dp := range ps {
		// TODO: this is temporally for debugging, getting every project
		// from s3 can be expensive
		v, err := s.GetProject(dp.ID)
		if err != nil {
			return []Project{}, err
		}
		p[i] = v
	}

	return p, nil
}

func (s *Service) UpdateProject(id string, project Project) error {
	s.assert.NotNil(s.db)
	s.assert.NotNil(s.s3)
	s.assert.NotZero(s.bucket)
	s.assert.NotNil(s.ctx)
	s.assert.NotZero(id)

	c, err := json.Marshal(project)
	if err != nil {
		return err
	}

	s.log.Debug("Updating project on storage", slog.String("id", id))

	f := fmt.Sprintf("%s.comic.json", id)
	_, err = s.s3.PutObject(s.ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Body:   bytes.NewReader(c),
		Key:    &f,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProject(id string) error {
	s.assert.NotNil(s.db)
	s.assert.NotNil(s.s3)
	s.assert.NotZero(s.bucket)
	s.assert.NotNil(s.ctx)
	s.assert.NotZero(id)

	err := s.db.DeleteProject(id)
	if err != nil {
		return err
	}

	f := fmt.Sprintf("%s.comic.json", id)
	_, err = s.s3.DeleteObject(s.ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &f,
	})
	if err != nil {
		return err
	}

	return nil
}
