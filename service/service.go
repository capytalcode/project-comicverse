package service

import (
	"context"
	"errors"
	"log/slog"

	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service struct {
	db     *database.Database
	s3     *s3.Client
	bucket string

	ctx context.Context

	assert tinyssert.Assertions
	log    *slog.Logger
}

func New(cfg Config) (*Service, error) {
	if cfg.DB == nil {
		return nil, errors.New("database should not be a nil pointer")
	}
	if cfg.S3 == nil {
		return nil, errors.New("s3 client should not be a nil pointer")
	}
	if cfg.Bucket == "" {
		return nil, errors.New("bucket should not be a empty string")
	}
	if cfg.Context == nil {
		return nil, errors.New("context should not be a nil interface")
	}
	if cfg.Assertions == nil {
		return nil, errors.New("assertions should not be a nil interface")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger should not be a nil pointer")
	}
	return &Service{
		db:     cfg.DB,
		s3:     cfg.S3,
		bucket: cfg.Bucket,

		ctx: cfg.Context,

		assert: cfg.Assertions,
		log:    cfg.Logger,
	}, nil
}

type Config struct {
	DB     *database.Database
	S3     *s3.Client
	Bucket string

	Context context.Context

	Assertions tinyssert.Assertions
	Logger     *slog.Logger
}
