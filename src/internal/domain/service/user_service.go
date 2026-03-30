package service

import (
	"context"
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/domain/model"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, req model.SignUpRequest) error
	GetUserBylogin(ctx context.Context, login string) (*model.User, error)
	GetUserInfo(ctx context.Context, login string) (*response.UserResponse, error)
	ValidateCredentials(ctx context.Context, login, password string) (uuid.UUID, error)
}
