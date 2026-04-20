package request

type MakeMoveRequest struct {
	GameID string `json:"game_id"`
	Row    int    `json:"row" validate:"required,min=0,max=2"`
	Col    int    `json:"col" validate:"required,min=0,max=2"`
	Player string `json:"player" validate:"required,oneof=X O"`
}
