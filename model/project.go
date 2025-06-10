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

func (m Project) Validate() error {
	errs := []error{}
	if len(m.ID) == 0 {
		errs = append(errs, ErrZeroValue{Name: "UUID"})
	}
	if m.Title == "" {
		errs = append(errs, ErrZeroValue{Name: "Title"})
	}
	if m.DateCreated.IsZero() {
		errs = append(errs, ErrZeroValue{Name: "DateCreated"})
	}
	if m.DateUpdated.IsZero() {
		errs = append(errs, ErrZeroValue{Name: "DateUpdated"})
	}

	if len(errs) > 0 {
		return ErrInvalidModel{Name: "Project", Errors: errs}
	}

	return nil
}
