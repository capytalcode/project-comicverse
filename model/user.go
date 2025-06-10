package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"` // Must be unique
	Password []byte    `json:"password"`

	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

func (u User) Validate() error {
	errs := []error{}
	if len(u.ID) == 0 {
		errs = append(errs, ErrZeroValue{Name: "ID"})
	}
	if u.Username == "" {
		errs = append(errs, ErrZeroValue{Name: "Username"})
	}
	if len(u.Password) == 0 {
		errs = append(errs, ErrZeroValue{Name: "Password"})
	}
	if u.DateCreated.IsZero() {
		errs = append(errs, ErrZeroValue{Name: "DateCreated"})
	}
	if u.DateUpdated.IsZero() {
		errs = append(errs, ErrZeroValue{Name: "DateUpdated"})
	}

	if len(errs) > 0 {
		return ErrInvalidModel{Name: "User", Errors: errs}
	}

	return nil
}
