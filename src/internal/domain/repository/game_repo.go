package db

import (
	"context"
	rp "tictactoe/internal/repository/model"
)

type GameRepository interface {
	FindGameById(ctx context.Context, id string) (*rp.GameModel, error)
	FindGameByWaitingStatus(ctx context.Context, id string) (*rp.GameModel, error)
	FindAllWaitingGames(ctx context.Context) ([]*rp.GameModel, error)
	GetGameHistoryByPlayerID(ctx context.Context, playerID string) ([]string, error)

	SaveGame(ctx context.Context, g rp.GameModel) error
	UpdateGame(ctx context.Context, g rp.GameModel) error
}
