package response

import (
	"time"
)

type GameResponse struct {
	GameID         string      `json:"id"`
	Board          [3][3]uint8 `json:"board"`
	FirstPlayerID  string      `json:"firstPlayer_ID"`
	SecondPlayerID string      `json:"secondPlayer_ID"`
	Status         string      `json:"status"`
	Winner         string      `json:"winner,omitempty"`
	CurrentTurn    string      `json:"current_turn"`
}

type CreateGameResponse struct {
	GameID    string    `json:"game_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type WaitingGamesResponse struct {
	Games []GameInfo `json:"games"`
	Count int        `json:"count"`
}

type GameInfo struct {
	GameID        string `json:"game_id"`
	FirstPlayerID string `json:"first_player_id"`
	Status        string `json:"status"`
}
