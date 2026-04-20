package model

import "time"

type UserModel struct {
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	UUID      string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type LeaderBoard struct {
	UUID    string  `db:"id"`
	Winrate float32 `db:"winrate"`
}
