package response

type UserResponse struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	GamesPlayed int    `json:"games_played"`
	Wins        int    `json:"wins"`
}
