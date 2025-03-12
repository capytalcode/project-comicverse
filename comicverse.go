package comicverse

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/database"
	"forge.capytal.company/capytalcode/project-comicverse/router"
	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/capytalcode/project-comicverse/static"
	"forge.capytal.company/capytalcode/project-comicverse/templates"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func New(cfg Config, opts ...Option) (http.Handler, error) {
	app := &app{
		db: cfg.DB,
		s3: cfg.S3,

		staticFiles:     static.Files(),
		developmentMode: false,
		ctx:             context.Background(),

		assert: tinyssert.NewAssertions(),
		logger: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	for _, opt := range opts {
		opt(app)
	}

	if app.db == nil {
		return nil, errors.New("database interface must not be nil")
	}
	if app.s3 == nil {
		return nil, errors.New("s3 client must not be nil")
	}

	if app.staticFiles == nil {
		return nil, errors.New("static files must not be a nil interface")
	}

	if app.ctx == nil {
		return nil, errors.New("context must not be a nil interface")
	}

	if app.logger == nil {
		return nil, errors.New("logger must not be a nil interface")
	}

	if app.assert == nil {
		return nil, errors.New("assertions must not be a nil interface")
	}

	return app, app.setup()
}

type Config struct {
	DB *sql.DB
	S3 *s3.Client
}

type Option func(*app)

func WithContext(ctx context.Context) Option {
	return func(app *app) { app.ctx = ctx }
}

func WithStaticFiles(f fs.FS) Option {
	return func(app *app) { app.staticFiles = f }
}

func WithAssertions(a tinyssert.Assertions) Option {
	return func(app *app) { app.assert = a }
}

func WithLogger(l *slog.Logger) Option {
	return func(app *app) { app.logger = l }
}

func WithDevelopmentMode() Option {
	return func(app *app) { app.developmentMode = true }
}

type app struct {
	db *sql.DB
	s3 *s3.Client

	ctx context.Context

	staticFiles     fs.FS
	developmentMode bool

	handler http.Handler

	assert tinyssert.Assertions
	logger *slog.Logger
}

func (app *app) setup() error {
	app.assert.NotNil(app.db)
	app.assert.NotNil(app.s3)
	app.assert.NotNil(app.ctx)
	app.assert.NotNil(app.staticFiles)
	app.assert.NotNil(app.logger)

	var err error

	database, err := database.New(database.Config{
		SQL:        app.db,
		Context:    app.ctx,
		Assertions: app.assert,
		Logger:     app.logger.WithGroup("database"),
	})
	if err != nil {
		return errors.Join(errors.New("unable to create database struct"), err)
	}

	service, err := service.New(service.Config{
		S3: app.s3,
		DB:     database,

		Context: app.ctx,

		Assertions: app.assert,
		Logger:     app.logger.WithGroup("service"),
	})
	if err != nil {
		return errors.Join(errors.New("unable to initiate service"), err)
	}

	app.handler, err = router.New(router.Config{
		Service: service,

		Templates:    templates.Templates(),
		DisableCache: app.developmentMode,
		StaticFiles:  app.staticFiles,

		Assertions: app.assert,
		Logger:     app.logger.WithGroup("router"),
	})
	if err != nil {
		return errors.Join(errors.New("unable to initiate router"), err)
	}

	return err
}

func (app *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.assert.NotNil(app.handler)
	app.handler.ServeHTTP(w, r)
}
