package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/capytalcode/project-comicverse/repository"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	repo *repository.User

	assert tinyssert.Assertions
	log    *slog.Logger
}

func NewUser(repo *repository.User, logger *slog.Logger, assert tinyssert.Assertions) *User {
	assert.NotNil(repo)
	assert.NotNil(logger)

	return &User{repo: repo, assert: assert, log: logger}
}

func (svc *User) Register(username, password string) (model.User, error) {
	svc.assert.NotNil(svc.repo)
	svc.assert.NotNil(svc.log)

	log := svc.log.With(slog.String("username", username))
	log.Info("Registering user")
	defer log.Info("Finished registering user")

	if _, err := svc.repo.GetByUsername(username); err == nil {
		return model.User{}, ErrUsernameAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, errors.New("service: unable to generate password hash")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return model.User{}, fmt.Errorf("service: unable to create user id", err)
	}

	now := time.Now()

	u := model.User{
		ID:          id,
		Username:    username,
		Password:    hash,
		DateCreated: now,
		DateUpdated: now,
	}

	u, err = svc.repo.Create(u)
	if err != nil {
		return model.User{}, fmt.Errorf("service: failed to create user model: %w", err)
	}

	return u, nil
}

func (svc *User) Login(username, password string) (user model.User, err error) {
	svc.assert.NotNil(svc.repo)
	svc.assert.NotNil(svc.log)

	log := svc.log.With(slog.String("username", username))
	log.Info("Logging in user")
	defer log.Info("Finished logging in user")

	user, err = svc.repo.GetByUsername(username)
	if err != nil {
		return model.User{}, fmt.Errorf("service: unable to find user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return model.User{}, fmt.Errorf("service: unable to compare passwords: %w", err)
	}

	return user, nil
}

var (
	ErrUsernameAlreadyExists = errors.New("service: username already exists")
	ErrPasswordTooLong       = bcrypt.ErrPasswordTooLong
	ErrIncorrectPassword     = bcrypt.ErrMismatchedHashAndPassword
)
