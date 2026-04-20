package mapper

import (
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/domain"
)

func DomainToResponse(game *domain.Game) *response.GameResponse {
	if game == nil {
		return nil
	}

	var winPlayer string
	var currentPlayer string

	switch {
	case game.Status == domain.StatusDraw:
		winPlayer = "draw"
	case game.Status == domain.StatusGameOver && game.Winner.ID != "":
		winPlayer = game.Winner.ID
	default:
		winPlayer = "nothing"
	}

	switch game.CurrentTurn.Icon {
	case 1:
		currentPlayer = "X"
	case 2:
		currentPlayer = "O"
	default:
		currentPlayer = "unknown"
	}

	return &response.GameResponse{
		GameID:         game.ID,
		Board:          game.Board.Cells,
		FirstPlayerID:  game.FirstPlayer.ID,
		SecondPlayerID: game.SecondPlayer.ID,
		Status:         string(game.Status),
		Winner:         winPlayer,
		CurrentTurn:    currentPlayer,
	}
}

func ResponseToDomain(resp *response.GameResponse) domain.Game {

	firstPlayer := domain.Player{
		ID:   resp.FirstPlayerID,
		Icon: 1,
	}
	secondPlayer := domain.Player{
		ID:   resp.SecondPlayerID,
		Icon: 2,
	}
	winner := domain.Player{}
	turn := domain.Player{}
	switch resp.Winner {
	case resp.FirstPlayerID:
		winner.ID = resp.FirstPlayerID
		winner.Icon = 1
	case resp.SecondPlayerID:
		winner.ID = resp.SecondPlayerID
		winner.Icon = 2
	}

	switch resp.CurrentTurn {
	case resp.FirstPlayerID:
		turn.ID = resp.FirstPlayerID
		turn.Icon = 1
	case resp.SecondPlayerID:
		turn.ID = resp.SecondPlayerID
		turn.Icon = 2
	}
	boardResp := struct {
		Cells [3][3]uint8
	}{Cells: resp.Board}
	return domain.Game{
		ID:           resp.GameID,
		Board:        boardResp,
		FirstPlayer:  firstPlayer,
		SecondPlayer: secondPlayer,
		Status:       resp.Status,
		Winner:       winner,
		CurrentTurn:  turn,
	}
}

func DomainToGameInfo(game *domain.Game) response.GameInfo {
	return response.GameInfo{
		GameID:        game.ID,
		FirstPlayerID: game.FirstPlayer.ID,
		Status:        string(game.Status),
	}
}

func DomainToWaitingGamesResponse(games []*domain.Game) response.WaitingGamesResponse {
	gameInfos := make([]response.GameInfo, 0, len(games))
	for _, game := range games {
		gameInfos = append(gameInfos, DomainToGameInfo(game))
	}

	return response.WaitingGamesResponse{
		Games: gameInfos,
		Count: len(gameInfos),
	}
}
