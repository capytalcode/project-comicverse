package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID // Must be unique, represented as base64 string in URLs
	Title       string    // Must not be empty
	DateCreated time.Time
	DateUpdated time.Time
}

var _ Model = (*Project)(nil)

func (p Project) Validate() error {
	errs := []error{}
	if len(p.ID) == 0 {
		errs = append(errs, ErrZeroValue{Name: "UUID"})
	}
	if p.Title == "" {
		errs = append(errs, ErrZeroValue{Name: "Title"})
	}
	if p.DateCreated.IsZero() {
		errs = append(errs, ErrZeroValue{Name: "DateCreated"})
	}
	if p.DateUpdated.IsZero() {
		errs = append(errs, ErrZeroValue{Name: "DateUpdated"})
	}

	if len(errs) > 0 {
		return ErrInvalidModel{Name: "Project", Errors: errs}
	}

	return nil
}
