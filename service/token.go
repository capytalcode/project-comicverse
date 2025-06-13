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
		UserID:      user.ID,
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

	_, ok := token.Claims.(jwt.RegisteredClaims) // TODO: Check issuer and if the token was issued at the correct date
	if !ok {
		return nil, errors.New("service: invalid claims type")
	}

	return token, nil
}

func (svc Token) Revoke(token *jwt.Token) error {
	svc.assert.NotNil(svc.log)
	svc.assert.NotNil(svc.repo)
	svc.assert.NotNil(token)

	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok {
		return errors.New("service: invalid claims type")
	}

	log := svc.log.With(slog.String("token_id", claims.ID))
	log.Info("Revoking token")
	defer log.Info("Finished revoking token")

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return errors.Join(errors.New("service: invalid token UUID"), err)
	}

	user, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errors.Join(errors.New("service: invalid token subject UUID"), err)
	}

	// TODO: Mark tokens as revoked instead of deleting them
	err = svc.repo.Delete(jti, user)
	if err != nil {
		return errors.Join(errors.New("service: failed to delete token"), err)
	}

	return nil
}

func (svc Token) IsRevoked(token *jwt.Token) (bool, error) {
	svc.assert.NotNil(svc.log)
	svc.assert.NotNil(svc.repo)
	svc.assert.NotNil(token)

	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok {
		return false, errors.New("service: invalid claims type")
	}

	log := svc.log.With(slog.String("token_id", claims.ID))
	log.Info("Checking if token is revoked")
	defer log.Info("Finished checking if token is revoked")

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return false, errors.Join(errors.New("service: invalid token UUID"), err)
	}

	user, err := uuid.Parse(claims.Subject)
	if err != nil {
		return false, errors.Join(errors.New("service: invalid token subject UUID"), err)
	}

	_, err = svc.repo.Get(jti, user)
	if errors.Is(err, repository.ErrNotFound) {
		return true, nil
	} else if err != nil {
		return false, errors.Join(errors.New("service: failed to get token"), err)
	}

	return false, nil
}
