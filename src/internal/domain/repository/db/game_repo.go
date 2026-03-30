package db

import (
	"context"
	rp "tictactoe/internal/repository/db/model"
)

type GameRepository interface {
	FindGameById(ctx context.Context, id string) (*rp.GameModel, error)
	FindGameByWaitingStatus(ctx context.Context, id string) (*rp.GameModel, error)
	FindAllWaitingGames(ctx context.Context) ([]*rp.GameModel, error)
	SaveGame(ctx context.Context, g rp.GameModel) error
	UpdateGame(ctx context.Context, g rp.GameModel) error
	RemoveGame(ctx context.Context, gameID string) error
}
