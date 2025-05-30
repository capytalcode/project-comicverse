package model

import (
	"time"
)

type User struct {
	Username string `json:"username"` // Must be unique
	Password []byte `json:"password"`

	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}
