package auth

import (
	"context"
	"errors"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/usecase/user"

	"github.com/google/uuid"
)

type AuthService struct {
	userService *user.UserService
}

type UserAuthenticator struct {
	Login    string
	Password string
}

func NewAuthService(userService *user.UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (a *AuthService) SignUp(ctx context.Context, req model.SignUpRequest) error {
	if err := validateCredentialsFormat(req.Login, req.Password); err != nil {
		return err
	}
	return a.userService.CreateUser(ctx, req)
}

func validateCredentialsFormat(login, password string) error {
	if login == "" {
		return errors.New("login cannot be empty")
	}
	if len(login) < 3 || len(login) > 50 {
		return errors.New("login must be between 3 and 50 characters")
	}

	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) < 6 || len(password) > 72 {
		return errors.New("password must be between 6 and 72 characters")
	}

	return nil
}

func (a *AuthService) SignIn(ctx context.Context, username, password string) (uuid.UUID, error) {
	if err := validateCredentialsFormat(username, password); err != nil {
		return uuid.Nil, err
	}
	return a.userService.ValidateCredentials(ctx, username, password)
}
