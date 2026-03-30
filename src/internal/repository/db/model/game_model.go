package model

import "time"

type GameModel struct {
	ID             string     `db:"id"`
	Board          BoardModel `db:"board"`
	FirstPlayerID  string     `db:"firstPlayer_id"`
	SecondPlayerID string     `db:"secondPlayer_id"`
	CurrentTurn    string     `db:"current_turn"`
	Status         string     `db:"status"`
	Winner         string     `db:"winner"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}
