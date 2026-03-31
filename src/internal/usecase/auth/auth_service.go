package auth

import (
	"context"
	"errors"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/usecase/user"
	"tictactoe/pkg/jwt"
)

type AuthService struct {
	userService *user.UserService
	jwtProvider *jwt.JwtProvider
}

type UserAuthenticator struct {
	Login    string
	Password string
}

func NewAuthService(userService *user.UserService, jwtProvider *jwt.JwtProvider) *AuthService {
	return &AuthService{
		userService: userService,
		jwtProvider: jwtProvider,
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

func (a *AuthService) SignIn(ctx context.Context, req jwt.JwtRequest) (*jwt.JwtResponse, error) {
	userID, err := a.userService.ValidateCredentials(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	user, err := a.userService.GetUserByID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	access, _ := a.jwtProvider.GenerateAccessToken(*user)
	refresh, _ := a.jwtProvider.GenerateRefreshToken(*user)

	return &jwt.JwtResponse{
		Type:         "Bearer",
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (a *AuthService) UpdateAccessToken(ctx context.Context, refreshToken string) (*jwt.JwtResponse, error) {
	if !a.jwtProvider.ValidateRefreshToken(refreshToken) {
		return nil, errors.New("invalid or expired refresh token")
	}

	userID, err := a.jwtProvider.GetUUIDByToken(refreshToken, true)
	if err != nil {
		return nil, err
	}
	user, err := a.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	newAccessToken, _ := a.jwtProvider.GenerateAccessToken(*user)
	newRefreshToken, _ := a.jwtProvider.GenerateRefreshToken(*user)

	return &jwt.JwtResponse{
		Type:         "Bearer",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *AuthService) UpdateRefreshToken(ctx context.Context, refreshToken string) (*jwt.JwtResponse, error) {
	return a.UpdateAccessToken(ctx, refreshToken)
}
