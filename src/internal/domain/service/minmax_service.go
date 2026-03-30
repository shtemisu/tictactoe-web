package service

import "tictactoe/internal/domain/model"

type MinMaxService interface {
	GetBestMove(game model.Game) (uint8, uint8)
}
