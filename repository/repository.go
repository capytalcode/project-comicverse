package repository

import (
	"errors"
	"time"

)
var (
	ErrDatabaseConn = errors.New("failed to begin transaction/connection with database")
	ErrExecuteQuery = errors.New("failed to execute query")
	ErrCommitQuery  = errors.New("failed to commit transaction")
	ErrInvalidData  = errors.New("data sent to save is invalid")
	ErrNotFound     = sql.ErrNoRows
)

var dateFormat = time.RFC3339
