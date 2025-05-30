package service

import (
	"time"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenService struct {
	assert tinyssert.Assertions
}

func NewTokenService(assert tinyssert.Assertions) *TokenService {
	return &TokenService{assert: assert}
}

func (s *TokenService) Issue(user model.User) (*jwt.Token, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		ID:        id.String(),
		Subject:   user.Username,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	})
}
