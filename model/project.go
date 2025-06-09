package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	UUID        uuid.UUID // Must be unique, represented as base64 string in URLs
	Title       string    // Must not be empty
	DateCreated time.Time
	DateUpdated time.Time
}

