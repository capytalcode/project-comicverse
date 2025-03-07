package comicverse

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/router"
	"forge.capytal.company/capytalcode/project-comicverse/static"
	"forge.capytal.company/loreddev/x/tinyssert"
)

func New(cfg Config, opts ...Option) (http.Handler, error) {
	app := &app{
		staticFiles:     static.Files(),
		developmentMode: false,
		context:         context.Background(),

		assert: tinyssert.NewAssertions(),
		logger: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	for _, opt := range opts {
		opt(app)
	}

	if app.staticFiles == nil {
		return nil, errors.New("static files must not be a nil interface")
	}

	if app.context == nil {
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
}

type Option func(*app)

func WithContext(ctx context.Context) Option {
	return func(app *app) { app.context = ctx }
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
	handler http.Handler

	staticFiles     fs.FS
	developmentMode bool
	context         context.Context

	assert tinyssert.Assertions
	logger *slog.Logger
}

func (app *app) setup() error {
	app.assert.NotNil(app.staticFiles)
	app.assert.NotNil(app.context)
	app.assert.NotNil(app.logger)

	var err error

	app.handler, err = router.New(router.Config{
		DisableCache: app.developmentMode,
		StaticFiles:  app.staticFiles,

		Assertions: app.assert,
		Logger:     app.logger,
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
