package domain

type MinMaxService interface {
	GetBestMove(game Game) (uint8, uint8)
}
