package model

import (
	"database/sql"
	"time"
)

type GameModel struct {
	ID             string
	Board          BoardModel
	FirstPlayerID  string
	SecondPlayerID sql.NullString // тоже может быть NULL
	CurrentTurn    sql.NullString // может быть NULL
	Status         string
	Winner         sql.NullString // МОЖЕТ БЫТЬ NULL!
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type EndedGames struct {
	ID     string
	Status string
}
