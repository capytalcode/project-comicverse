package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"

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

type ProjectPage struct {
	ID           string                     `json:"id"`
	Interactions map[string]PageInteraction `json:"interactions"`
	Image        io.ReadCloser              `json:"-"`
}

type PageInteraction struct {
	URL string `json:"url"`
	X   uint16 `json:"x"`
	Y   uint16 `json:"y"`
}

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

func (s *Service) GetPage(projectID string, pageID string) (ProjectPage, error) {
	s.assert.NotNil(s.ctx)
	s.assert.NotNil(s.s3)
	s.assert.NotNil(s.bucket)
	s.assert.NotZero(projectID)
	s.assert.NotNil(pageID)

	p, err := s.GetProject(projectID)
	if err != nil {
		return ProjectPage{}, errors.Join(errors.New("unable to get project"), err)
	}

	pageIndex := slices.IndexFunc(p.Pages, func(p ProjectPage) bool { return p.ID == pageID })
	if pageIndex == -1 {
		return ProjectPage{}, ErrPageNotExists
	}
	page := p.Pages[pageIndex]

	k := fmt.Sprintf("%s/%s", projectID, pageID)
	res, err := s.s3.GetObject(s.ctx, &s3.GetObjectInput{
		Key:    &k,
		Bucket: &s.bucket,
	})
	if err != nil {
		var resErr *awshttp.ResponseError
		if errors.As(err, &resErr) && resErr.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			// TODO: This would probably be better in some background "maintenance" worker
			p.Pages = slices.Delete(p.Pages, pageIndex, pageIndex)
			_ = s.UpdateProject(projectID, p)

			return ProjectPage{}, errors.Join(ErrPageNotExists, resErr)
		}
		return ProjectPage{}, err
	}

	s.assert.NotNil(res.Body)
	s.assert.NotNil(page.Interactions)

	page.Image = res.Body

	return page, nil
}

func (s *Service) UpdatePage(projectID string, page ProjectPage) error {
	s.assert.NotNil(s.ctx)
	s.assert.NotNil(s.s3)
	s.assert.NotNil(s.bucket)
	s.assert.NotZero(projectID)
	s.assert.NotZero(page.ID)
	s.assert.NotNil(page.Interactions)

	p, err := s.GetProject(projectID)
	if err != nil {
		return errors.Join(errors.New("unable to get project"), err)
	}

	pageIndex := slices.IndexFunc(p.Pages, func(p ProjectPage) bool { return p.ID == page.ID })
	if pageIndex == -1 {
		return ErrPageNotExists
	}
	p.Pages[pageIndex] = page

	// TODO: Probably a "lastUpdated" timestamp in the ProjectPage data would be better
	// so we don't update equal images. Changing the image in ProjectPage would be better
	// using a method, or could be completely decoupled from the struct.
	if page.Image != nil {
		k := fmt.Sprintf("%s/%s", projectID, page.ID)
		_, err = s.s3.PutObject(s.ctx, &s3.PutObjectInput{
			Key:    &k,
			Body:   page.Image,
			Bucket: &s.bucket,
		})
		if err != nil {
			return errors.Join(errors.New("error while trying to update image"), err)
		}
	}

	err = s.UpdateProject(projectID, p)
	if err != nil {
		return errors.Join(errors.New("error while trying to update project"), err)
	}

	return nil
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
