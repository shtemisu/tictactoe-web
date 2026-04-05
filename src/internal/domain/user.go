package domain

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Login    string
	Password string
}

type UserResponse struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	GamesPlayed int    `json:"games_played"`
	Wins        int    `json:"wins"`
}

type LeaderBoardResponse struct {
	UUID    string  `json:"id"`
	Winrate float32 `json:"winrate"`
}

type UserService interface {
	CreateUser(ctx context.Context, req SignUpRequest) error
	GetUserBylogin(ctx context.Context, login string) (*User, error)
	GetUserInfo(ctx context.Context, ID string) (*UserResponse, error)
	GetLeaderBoard(ctx context.Context, limit string) ([]LeaderBoardResponse, error)
	ValidateCredentials(ctx context.Context, login, password string) (uuid.UUID, error)
}

type SignUpRequest struct {
	Login    string
	Password string
}
