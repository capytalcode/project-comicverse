package model

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	DateCreated time.Time
	DateExpires time.Time
}

func (t Token) Validate() error {
	errs := []error{}
	if len(t.ID) == 0 {
		errs = append(errs, ErrZeroValue{Name: "ID"})
	}
	if len(t.UserID) == 0 {
		errs = append(errs, ErrZeroValue{Name: "User"})
	}
	if t.DateCreated.IsZero() {
		errs = append(errs, ErrZeroValue{Name: "DateCreated"})
	}
	if t.DateExpires.IsZero() {
		errs = append(errs, ErrZeroValue{Name: "DateExpires"})
	}
	if len(errs) > 0 {
		return ErrInvalidModel{Name: "Token", Errors: errs}
	}
	return nil
}
