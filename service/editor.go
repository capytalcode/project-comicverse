package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/internals/randstr"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const pageIDLength = 6

var ErrPageNotExists = errors.New("page does not exists in storage")

type Project struct {
	ID    string        `json:"id"`
	Title string        `json:"title"`
	Pages []ProjectPage `json:"pages"`
}

type ProjectPage struct{}

func (s *Service) AddPage(projectID string, img io.Reader) error {
	s.assert.NotNil(s.ctx)
	s.assert.NotNil(s.s3)
	s.assert.NotNil(s.bucket)
	s.assert.NotZero(projectID)
	s.assert.NotNil(img)

	id, err := randstr.NewHex(pageIDLength)
	if err != nil {
		return err
	}

	p, err := s.GetProject(projectID)
	if err != nil {
		return errors.Join(errors.New("unable to get project"), err)
	}

	p.Pages = append(p.Pages, ProjectPage{ID: id, Interactions: map[string]PageInteraction{}})

	k := fmt.Sprintf("%s/%s", projectID, id)
	_, err = s.s3.PutObject(s.ctx, &s3.PutObjectInput{
		Key:    &k,
		Body:   img,
		Bucket: &s.bucket,
	})
	if err != nil {
		return err
	}

	err = s.UpdateProject(projectID, p)
	return err
}

func (s *Service) GetPage(projectID string, imgID string) (io.Reader, error) {
	s.assert.NotNil(s.ctx)
	s.assert.NotNil(s.s3)
	s.assert.NotNil(s.bucket)
	s.assert.NotZero(projectID)
	s.assert.NotNil(imgID)

	k := fmt.Sprintf("%s/%s", projectID, imgID)
	res, err := s.s3.GetObject(s.ctx, &s3.GetObjectInput{
		Key:    &k,
		Bucket: &s.bucket,
	})
	if err != nil {
		var resErr *awshttp.ResponseError
		if errors.As(err, &resErr) && resErr.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return nil, errors.Join(ErrPageNotExists, resErr)
		}
		return nil, err
	}

	s.assert.NotNil(res.Body)
	return res.Body, nil
}

func (s *Service) DeletePage(projectID string, id string) error {
	s.assert.NotNil(s.ctx)
	s.assert.NotNil(s.s3)
	s.assert.NotNil(s.bucket)
	s.assert.NotZero(projectID)
	s.assert.NotNil(id)

	p, err := s.GetProject(projectID)
	if err != nil {
		return errors.Join(errors.New("unable to get project"), err)
	}

	k := fmt.Sprintf("%s/%s", projectID, id)
	_, err = s.s3.DeleteObject(s.ctx, &s3.DeleteObjectInput{
		Key:    &k,
		Bucket: &s.bucket,
	})
	if err != nil {
		return err
	}

	p.Pages = slices.DeleteFunc(p.Pages, func(p ProjectPage) bool { return p.ID == id })

	err = s.UpdateProject(projectID, p)
	return err
}
