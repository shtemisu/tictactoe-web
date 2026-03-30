package service

import (
	"context"
	"tictactoe/internal/domain/model"
)

type GameService interface {
	CreateGameWithAI(ctx context.Context, playerID string) (string, error)
	CreateMultiplayerGame(ctx context.Context, firstPlayerID string) (string, error)
	JoinToGame(ctx context.Context, secondPlayerID string, gameID string) (*model.Game, error)
	GetGame(ctx context.Context, gameID string) (*model.Game, error)
	GetAllWaitingGames(ctx context.Context) ([]*model.Game, error)
	GetTurn(ctx context.Context, gameID string) (uint8, error)
	DoMove(ctx context.Context, gameID string, row uint8, col uint8) (bool, error)
	GetNextMove(ctx context.Context, gameID string) (uint8, uint8, error)
}
