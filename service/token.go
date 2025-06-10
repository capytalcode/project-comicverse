package service

import (
	"crypto/ed25519"
	"errors"
	"log/slog"
	"time"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/capytalcode/project-comicverse/repository"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Token struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey

	repo *repository.Token

	log    *slog.Logger
	assert tinyssert.Assertions
}

func NewToken(
	privateKey ed25519.PrivateKey,
	publicKey ed25519.PublicKey,
	repo *repository.Token,
	logger *slog.Logger,
	assert tinyssert.Assertions,
) *Token {
	assert.NotZero(privateKey)
	assert.NotZero(publicKey)
	assert.NotZero(repo)
	assert.NotZero(logger)

	return &Token{assert: assert}
}

func (svc *Token) Issue(user model.User) (string, error) { // TODO: Return a refresh token
	svc.assert.NotNil(svc.privateKey)
	svc.assert.NotNil(svc.log)
	svc.assert.NotZero(user)

	log := svc.log.With(slog.String("user_id", user.ID.String()))
	log.Info("Issuing new token")
	defer log.Info("Finished issuing token")

	jti, err := uuid.NewV7()
	if err != nil {
		return "", errors.Join(errors.New("service: failed to generate token UUID"), err)
	}

	now := time.Now()
	expires := now.Add(30 * 24 * time.Hour) // TODO: Make the JWT short lived and use refresh tokens to create new JWTs

	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{
		Issuer:    "comicverse", // TODO: Make application ID and Name be a parameter
		Subject:   user.ID.String(),
		Audience:  jwt.ClaimStrings{"comicverse"}, // TODO: When we have third-party apps integration, this should be the name/URI/id of the app
		ExpiresAt: jwt.NewNumericDate(expires),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        jti.String(),
	})

	signed, err := t.SignedString(svc.privateKey)
	if err != nil {
		return "", errors.Join(errors.New("service: failed to sign token"), err)
	}

	// TODO: Store refresh tokens in repo
	err = svc.repo.Create(model.Token{
		ID:          jti,
		DateCreated: now,
		DateExpires: expires,
	})
	if err != nil {
		return "", errors.Join(errors.New("service: failed to save token"), err)
	}

	return signed, nil
}

func (svc Token) Parse(tokenStr string) (*jwt.Token, error) {
	svc.assert.NotNil(svc.publicKey)

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return svc.publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}))
	if err != nil {
		return nil, errors.Join(errors.New("service: invalid token"), err)
	}

	_, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("service: invalid claims type")
	}

	return token, nil
}

