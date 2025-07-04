package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"forge.capytal.company/loreddev/x/tinyssert"
)

type RoleRepository struct {
	db *sql.DB

	ctx    context.Context
	log    *slog.Logger
	assert tinyssert.Assertions
}
