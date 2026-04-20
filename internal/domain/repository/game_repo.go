package db

import (
	"context"
	"tictactoe/internal/repository/model"
	rp "tictactoe/internal/repository/model"
)

type GameRepository interface {
	FindGameById(ctx context.Context, id string) (*rp.GameModel, error)
	FindGameByWaitingStatus(ctx context.Context, id string) (*rp.GameModel, error)
	FindAllWaitingGames(ctx context.Context) ([]*rp.GameModel, error)
	GetGameHistoryByPlayerID(ctx context.Context, playerID string) ([]model.EndedGames, error)

	SaveGame(ctx context.Context, g rp.GameModel) error
	UpdateGame(ctx context.Context, g rp.GameModel) error
}
