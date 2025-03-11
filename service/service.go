package service

import (
	"database/sql"
	"errors"
	"log/slog"

	"forge.capytal.company/loreddev/x/tinyssert"
)

type service struct {
	assert tinyssert.Assertions
	log    *slog.Logger
}

func New(cfg Config) (Service, error) {
	if cfg.DB == nil {
		return nil, errors.New("database should not be a nil interface")
	}
	if cfg.S3 == nil {
		return nil, errors.New("s3 client should not be a nil interface")
	}
	if cfg.Assertions == nil {
		return nil, errors.New("assertions should not be a nil interface")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger should not be a nil pointer")
	}
	return &service{
		assert: cfg.Assertions,
		log:    cfg.Logger,
	}, nil
}

type Config struct {
	Assertions tinyssert.Assertions
	Logger     *slog.Logger
}

type Service interface {
	ListProjects()
	NewProject()
}
