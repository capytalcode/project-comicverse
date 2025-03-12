package database

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"forge.capytal.company/loreddev/x/tinyssert"
)

var ErrNoRows = sql.ErrNoRows

type Database struct {
	sql *sql.DB
	ctx context.Context

	assert tinyssert.Assertions
	log    *slog.Logger
}

func New(cfg Config) (*Database, error) {
	if cfg.SQL == nil {
		return nil, errors.New("SQL database interface should not be nil")
	}
	if cfg.Context == nil {
		return nil, errors.New("context interface should not be nil")
	}
	if cfg.Assertions == nil {
		return nil, errors.New("assertions interface should not be nil")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger should not be a nil pointer")
	}

	db := &Database{
		sql: cfg.SQL,
		ctx: cfg.Context,

		assert: cfg.Assertions,
		log:    cfg.Logger,
	}

	if err := db.setup(); err != nil {
		return nil, errors.Join(errors.New("error while setting up Database struct"), err)
	}

	return db, nil
}

type Config struct {
	SQL     *sql.DB
	Context context.Context

	Assertions tinyssert.Assertions
	Logger     *slog.Logger
}

func (db *Database) setup() error {
	db.assert.NotNil(db.sql)
	db.assert.NotNil(db.ctx)
	db.assert.NotNil(db.log)

	log := db.log
	log.Info("Setting up database")

	log.Debug("Pinging database")

	err := db.sql.PingContext(db.ctx)
	if err != nil {
		return errors.Join(errors.New("unable to ping database"), err)
	}

	log.Debug("Creating tables")

	tx, err := db.sql.BeginTx(db.ctx, nil)
	if err != nil {
		return errors.Join(errors.New("unable to start transaction to create tables"), err)
	}

	setups := []func(*sql.Tx) error{
		db.setupProjects,
	}

	for _, setup := range setups {
		err := setup(tx)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("unable to run transaction to create tables"), err)
	}

	return nil
}
