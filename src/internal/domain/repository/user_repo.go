package db

import (
	"context"
	rp "tictactoe/internal/repository/db/model"
)

type UserRepository interface {
	SaveUser(ctx context.Context, g rp.UserModel) error
	FindUserByLogin(ctx context.Context, login string) (*rp.UserModel, error)
	FindUserByID(ctx context.Context, ID string) (*rp.UserModel, error)
	GetUserStats(ctx context.Context, ID string) (gamesPlayed int, wins int, err error)
	RemoveUser(ctx context.Context, user_id string) error
}
