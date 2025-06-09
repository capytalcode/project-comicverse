package model

import (
	"fmt"
)

type Model interface {
	Validate() error
}

type ErrInvalidModel struct {
	Name   string
	Errors []error
}

var _ error = ErrInvalidModel{}

func (err ErrInvalidModel) Error() string {
	return fmt.Sprintf("model %q is invalid", err.Name)
}

type ErrInvalidValue struct {
	Name     string
	Actual   any
	Expected []any
}

var _ error = ErrInvalidValue{}

func (err ErrInvalidValue) Error() string {
	var msg string

	if err.Name != "" {
		msg = fmt.Sprintf("%q has", err.Name)
	}

	msg = msg + " incorrect value"

	if err.Actual != nil {
		msg = msg + fmt.Sprintf(" %q", err.Actual)
	}

	if len(err.Expected) == 0 || err.Expected == nil {
		return msg
	}

	msg = fmt.Sprintf("%s, expected %q", msg, err.Expected[0])
	if len(err.Expected) > 1 {
		if len(err.Expected) == 2 {
			msg = msg + fmt.Sprintf(" or %q", err.Expected[1])
		} else {
			for v := range err.Expected[1 : len(err.Expected)-1] {
				msg = msg + fmt.Sprintf(", %q", v)
			}
			msg = msg + fmt.Sprintf(", or %q", err.Expected[len(err.Expected)-1])
		}
	}

	return msg
}

type ErrZeroValue ErrInvalidValue

func (err ErrZeroValue) Error() string {
	return fmt.Sprintf("%q has incorrect value, expected non-zero/non-empty value", err.Name)
}
