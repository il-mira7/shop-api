package service

import (
	"context"
	"errors"
	"fmt"
	"shopAPI/internal/authserver/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash []byte) (bool, error)
}

type AuthService struct {
	repo UserRepository
	jwt  *JwtProvider
}

func NewAuthService(repo UserRepository, jwt *JwtProvider) *AuthService {
	return &AuthService{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.User, password string) (string, error) {

	if password == "" {
		return "", &domain.EmptyPayloadError{Entity: "password"}
	}

	if user.Email == "" {
		return "", &domain.EmptyPayloadError{Entity: "email"}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashed)

	err = s.repo.CreateUser(ctx, &user)
	if err != nil {
		return "", err
	}

	return s.jwt.GenerateAccessToken(user.ID)
}

func (s *AuthService) Authenticate(ctx context.Context, email, password string) (string, error) {
	id, err := s.check(ctx, email, password)
	if err != nil {
		return "", err
	}
	return s.jwt.GenerateAccessToken(id)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", &domain.EmptyPayloadError{Entity: "token"}
	}

	userID, err := s.jwt.ValidateAccessToken(token)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}
	return userID, nil
}
func (s *AuthService) ChangePassword(ctx context.Context, email, passwordLast, passwordNew string) (bool, error) {
	if passwordNew == "" {
		return false, &domain.EmptyPayloadError{Entity: "passwordNew"}
	}
	id, err := s.check(ctx, email, passwordLast)
	if err != nil {
		return false, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwordNew), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("failed to hash password: %w", err)
	}
	return s.repo.UpdatePasswordHash(ctx, id, hashed)
}
func (s *AuthService) ResetPassword(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, &domain.EmptyPayloadError{Entity: "email"}
	}
	user, err := s.repo.GetByLogin(ctx, email)
	if err != nil {
		return false, err
	}
	password := generatePassword()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("failed to hash password: %w", err)
	}
	if _, err = s.repo.UpdatePasswordHash(ctx, user.ID, hashed); err != nil {
		return false, err
	}
	fmt.Println("──────────────────────────────────────────────────────────────────────────────────────────────")
	fmt.Println("")
	fmt.Printf("  Новый пароль для %s: %s\n", user.Email, password)
	fmt.Println("")
	fmt.Println("──────────────────────────────────────────────────────────────────────────────────────────────")
	return true, nil
}

func (s *AuthService) check(ctx context.Context, email, password string) (uuid.UUID, error) {
	if password == "" {
		return uuid.Nil, &domain.EmptyPayloadError{Entity: "password"}
	}
	if email == "" {
		return uuid.Nil, &domain.EmptyPayloadError{Entity: "email"}
	}

	user, err := s.repo.GetByLogin(ctx, email)
	if err != nil {
		var notFound *domain.NotFoundError
		if errors.As(err, &notFound) {
			return uuid.Nil, domain.ErrInvalidCredentials
		}
		return uuid.Nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return uuid.Nil, domain.ErrInvalidCredentials
	}
	return user.ID, nil
}

func generatePassword() string {
	password := uuid.NewString()
	return password[0:8]
}
