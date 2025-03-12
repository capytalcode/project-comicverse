package service

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"forge.capytal.company/capytalcode/project-comicverse/database"
	"forge.capytal.company/capytalcode/project-comicverse/internals/randstr"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const projectIDLength = 6

var (
	ErrProjectNotExists   = errors.New("project does not exists in database")
	ErrProjectInvalidUUID = errors.New("UUID provided is invalid")
)

type Project struct {
	XMLName  xml.Name `xml:"body"`
	ID       string   `xml:"id,attr"`
	Title    string   `xml:"h1"`
	Contents string   `xml:"-"`
}

func (s *Service) CreateProject() (Project, error) {
	s.assert.NotNil(s.db)
	s.assert.NotNil(s.s3)
	s.assert.NotNil(s.ctx)
	s.assert.NotZero(s.bucket)

	s.log.Debug("Creating new project")

	id, err := randstr.NewHex(projectIDLength)
	if err != nil {
		return Project{}, err
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
	}

	c, err := xml.Marshal(p)
	if err != nil {
		return Project{}, err
	}

	s.log.Debug("Creating project on storage", slog.String("id", id))

	f := fmt.Sprintf("%s.comic.xml", id)
	_, err = s.s3.PutObject(s.ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &f,
		Body:   bytes.NewReader(c),
	})
	if err != nil {
		return Project{}, err
	}

	p.Contents = string(c)

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

	p := Project{
		ID:    res.ID,
		Title: res.Title,
	}

	f := fmt.Sprintf("%s.comic.xml", p.ID)
	file, err := s.s3.GetObject(s.ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &f,
	})
	if err != nil {
		return p, err
	}

	c, err := io.ReadAll(file.Body)
	if err != nil {
		return p, err
	}

	p.Contents = string(c)

	return p, nil
}
