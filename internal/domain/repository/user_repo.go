package db

import (
	"context"
	"tictactoe/internal/domain"
	"tictactoe/internal/repository/model"
	rp "tictactoe/internal/repository/model"
)

type UserRepository interface {
	SaveUser(ctx context.Context, g rp.UserModel) error
	FindUserByLogin(ctx context.Context, login string) (*rp.UserModel, error)
	FindUserByID(ctx context.Context, ID string) (*rp.UserModel, error)
	GetUserStats(ctx context.Context, ID string) (gamesPlayed int, wins int, err error)
	GetLeaderBoard(ctx context.Context, limit string) ([]domain.LeaderBoardResponse, error)
	GetGameHistoryByPlayerID(ctx context.Context, playerID string) ([]model.EndedGames, error)
	RemoveUser(ctx context.Context, user_id string) error
}
