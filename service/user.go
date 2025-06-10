package service

import (
	"errors"
	"time"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/capytalcode/project-comicverse/repository"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	assert tinyssert.Assertions
	repo   *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository, assert tinyssert.Assertions) (*UserService, error) {
	if err := assert.NotNil(repo); err != nil {
		return nil, err
	}

	return &UserService{repo: repo, assert: assert}, nil
}

func (s *UserService) Register(username, password string) (model.User, error) {
	s.assert.NotNil(s.repo)

	if _, err := s.repo.GetByUsername(username); err == nil {
		return model.User{}, ErrAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	u := model.User{
		Username:    username,
		Password:    hash,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	u, err = s.repo.Create(u)
	if err != nil {
		return model.User{}, errors.Join(errors.New("failed to create user model"), err)
	}

	return u, nil
}

func (s *UserService) Login(username, password string) (signedToken string, user model.User, err error) {
	s.assert.NotNil(s.repo)

	user, err = s.repo.GetByUsername(username)
	if err != nil {
		return "", model.User{}, errors.Join(errors.New("unable to find user"), err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return "", model.User{}, errors.Join(errors.New("unable to compare passwords"), err)
	}

		return "", user, errors.Join(errors.New("unable to sign token"), err)
		return model.User{}, errors.Join(errors.New("service: unable to compare passwords"), err)
	}

	return user, nil
}

var (
	ErrAlreadyExists     = errors.New("model already exists")
	ErrNotFound          = repository.ErrNotFound
	ErrPasswordTooLong   = bcrypt.ErrPasswordTooLong
	ErrIncorrectPassword = bcrypt.ErrMismatchedHashAndPassword
)
