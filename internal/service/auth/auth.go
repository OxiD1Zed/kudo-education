package auth_service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"kode-education/internal/domain"
	"kode-education/internal/storage"
	"log/slog"
	"time"
)

type AuthService struct {
	log             *slog.Logger
	userProvider    UserProvider
	sessionProvider SessionProvider
}

type UserProvider interface {
	GetByUsername(username string) (domain.User, error)
}

type SessionProvider interface {
	Save(token string, idUser uint)
	GetByToken(token string) (uint, error)
}

func New(
	log *slog.Logger,
	userProvider UserProvider,
	sessionProvider SessionProvider,
) *AuthService {
	return &AuthService{
		log,
		userProvider,
		sessionProvider,
	}
}

func (a *AuthService) Authenticate(username string, password string) (string, error) {
	const op = "AuthService.Authenticate"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("attempting to authenticate user")

	user, err := a.userProvider.GetByUsername(username)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		a.log.Error("failed to get user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if checkPasswordHash(password, user.Password) {
		log.Info("user authenticate in successfully")

		token := generateToken(user)
		a.sessionProvider.Save(token, user.Id)

		return token, nil
	}

	a.log.Info("invalid credentials")

	return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
}

func (a *AuthService) Authorization(token string) (uint, error) {
	const op = "AuthService.Authorization"

	log := a.log.With(slog.String("op", op))

	log.Info("attempting to login user")

	idUser, err := a.sessionProvider.GetByToken(token)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			a.log.Warn("token not found", slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s: %w", op, ErrorInvalidToken)
		}

		a.log.Warn("failed to get token", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("the user has successfully logged in")

	return idUser, nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func checkPasswordHash(password, hash string) bool {
	return hashPassword(password) == hash
}

func generateToken(user domain.User) string {
	return hex.EncodeToString([]byte(time.Now().String() + user.Username + user.Password))
}
