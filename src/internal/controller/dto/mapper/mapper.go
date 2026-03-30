package mapper

import (
	"tictactoe/internal/controller/dto/response"
	domainModel "tictactoe/internal/domain/model"
)

func DomainToResponse(game *domainModel.Game) *response.GameResponse {
	if game == nil {
		return nil
	}

	var winPlayer string
	var currentPlayer string

	switch {
	case game.Status == domainModel.StatusDraw:
		winPlayer = "draw"
	case game.Status == domainModel.StatusGameOver && game.Winner.ID != "":
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
		Board:          game.Board,
		FirstPlayerID:  game.FirstPlayer.ID,
		SecondPlayerID: game.SecondPlayer.ID,
		Status:         string(game.Status),
		Winner:         winPlayer,
		CurrentTurn:    currentPlayer,
	}
}

func ResponseToDomain(resp *response.GameResponse) domainModel.Game {

	firstPlayer := domainModel.Player{
		ID:   resp.FirstPlayerID,
		Icon: 1,
	}
	secondPlayer := domainModel.Player{
		ID:   resp.SecondPlayerID,
		Icon: 2,
	}
	winner := domainModel.Player{}
	turn := domainModel.Player{}
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

	return domainModel.Game{
		ID:           resp.GameID,
		Board:        resp.Board,
		FirstPlayer:  firstPlayer,
		SecondPlayer: secondPlayer,
		Status:       resp.Status,
		Winner:       winner,
		CurrentTurn:  turn,
	}
}

func DomainToGameInfo(game *domainModel.Game) response.GameInfo {
	return response.GameInfo{
		GameID:        game.ID,
		FirstPlayerID: game.FirstPlayer.ID,
		Status:        string(game.Status),
	}
}

func DomainToWaitingGamesResponse(games []*domainModel.Game) response.WaitingGamesResponse {
	gameInfos := make([]response.GameInfo, 0, len(games))
	for _, game := range games {
		gameInfos = append(gameInfos, DomainToGameInfo(game))
	}

	return response.WaitingGamesResponse{
		Games: gameInfos,
		Count: len(gameInfos),
	}
}
