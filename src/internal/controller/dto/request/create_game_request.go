package request

type CreateGameRequest struct {
	Firstplayer string `json:"firstplayer" validate:"required,oneof=X O"`
}
