package service

import (
	"crypto/ed25519"
	"errors"
	"fmt"
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

func NewToken(cfg TokenConfig) *Token {
	cfg.Assertions.NotZero(cfg.PrivateKey)
	cfg.Assertions.NotZero(cfg.PublicKey)
	cfg.Assertions.NotZero(cfg.Repository)
	cfg.Assertions.NotZero(cfg.Logger)

	return &Token{
		privateKey: cfg.PrivateKey,
		publicKey:  cfg.PublicKey,
		repo:       cfg.Repository,
		log:        cfg.Logger,
		assert:     cfg.Assertions,
	}
}

type TokenConfig struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Repository *repository.Token
	Logger     *slog.Logger
	Assertions tinyssert.Assertions
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
		return "", fmt.Errorf("service: failed to generate token UUID: %w", err)
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
		return "", fmt.Errorf("service: failed to sign token: %w", err)
	}

	// TODO: Store refresh tokens in repo
	err = svc.repo.Create(model.Token{
		ID:          jti,
		UserID:      user.ID,
		DateCreated: now,
		DateExpires: expires,
	})
	if err != nil {
		return "", fmt.Errorf("service: failed to save token: %w", err)
	}

	return signed, nil
}

func (svc Token) Parse(tokenStr string) (*jwt.Token, error) {
	svc.assert.NotNil(svc.publicKey)
	svc.assert.NotNil(svc.log)

	log := svc.log.With(slog.String("preview_token", tokenStr[0:5]))
	log.Info("Parsing token")
	defer log.Info("Finished parsing token")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return svc.publicKey, nil
	}, jwt.WithValidMethods([]string{(&jwt.SigningMethodEd25519{}).Alg()}))
	if err != nil {
		log.Error("Invalid token", slog.String("error", err.Error()))
		return nil, fmt.Errorf("service: invalid token: %w", err)
	}

	// TODO: Check issuer and if the token was issued at the correct date
	// TODO: Structure token claims type
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Error("Invalid claims type", slog.String("claims", fmt.Sprintf("%#v", token.Claims)))
		return nil, fmt.Errorf("service: invalid claims type")
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
		return fmt.Errorf("service: invalid token UUID: %w", err)
	}

	user, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fmt.Errorf("service: invalid token subject UUID: %w", err)
	}

	// TODO: Mark tokens as revoked instead of deleting them
	err = svc.repo.Delete(jti, user)
	if err != nil {
		return fmt.Errorf("service: failed to delete token: %w", err)
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
		return false, fmt.Errorf("service: invalid token UUID: %w", err)
	}

	user, err := uuid.Parse(claims.Subject)
	if err != nil {
		return false, fmt.Errorf("service: invalid token subject UUID: %w", err)
	}

	_, err = svc.repo.Get(jti, user)
	if errors.Is(err, repository.ErrNotFound) {
		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("service: failed to get token: %w", err)
	}

	return false, nil
}
