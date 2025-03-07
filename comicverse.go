package comicverse

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"

	"forge.capytal.company/loreddev/x/tinyssert"
)

func New(cfg Config, opts ...Option) (http.Handler, error) {
	app := &app{
		developmentMode: false,
		context:         context.Background(),

		assert: tinyssert.NewAssertions(),
		logger: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	for _, opt := range opts {
		opt(app)
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
	developmentMode bool
	context         context.Context

	assert tinyssert.Assertions
	logger *slog.Logger
}

func (app *app) setup() error {
	app.assert.NotNil(app.context)
	app.assert.NotNil(app.logger)

	var err error
	return err
}
func (app *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
