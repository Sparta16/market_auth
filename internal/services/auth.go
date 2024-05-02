package services

import (
	"context"
	"errors"
	"fmt"
	uuid3 "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
	models "weather_back_api_getway/internal"
	"weather_back_api_getway/internal/repositories"
	"weather_back_api_getway/internal/utilites"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Login(username, email, password string) error
	Register(username, password string) error
}

type UserSaver interface {
	SaveUser(ctx context.Context, uuid string, username string, password []byte, email string) (string, error)
}

type UserProvider interface {
	User(
		ctx context.Context,
		login string,
	) (models.User, error)
}

type AuthServiceImpl struct {
	Database     *repositories.Database
	UserSaver    UserSaver
	UserProvider UserProvider
	TokenTTL     time.Duration
}

func New(d *repositories.Database, saver UserSaver, provider UserProvider, ttl time.Duration) *AuthServiceImpl {
	return &AuthServiceImpl{
		Database:     d,
		UserSaver:    saver,
		UserProvider: provider,
		TokenTTL:     ttl,
	}
}

func (a *AuthServiceImpl) Register(ctx context.Context, username, email, password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Info("failed to hash password", err)
		return "", err
	}

	uniqueID := uuid3.New()
	uuid := uniqueID.String()

	res, err := a.UserSaver.SaveUser(ctx, uuid, username, passHash, email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserExists) {
			return "", fmt.Errorf("user %s already exists", repositories.ErrUserExists)
		}
		slog.Info("failed to save user", err)
		return "", err
	}
	return res, nil
}

func (a *AuthServiceImpl) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.UserProvider.User(ctx, username)
	if err != nil {
		if errors.Is(err, repositories.ErrUserExists) {
			slog.Info("user exists", username)
			return "", fmt.Errorf("%s: %w", username, repositories.ErrUserExists)
		}
		slog.Info("failed to get user", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		slog.Info("invalid credentials", username)
		return "", fmt.Errorf("%s: %w", username, ErrInvalidCredentials)
	}

	token, err := utilites.NewToken(user, a.TokenTTL)
	if err != nil {
		slog.Info("failed to generate token", err)
		return "", err
	}
	return token, nil
}
