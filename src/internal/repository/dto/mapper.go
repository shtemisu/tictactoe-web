package dto

import (
	"errors"
	"log"
	domainModel "tictactoe/internal/domain/model"
	rp "tictactoe/internal/repository/model"
	"time"

	"github.com/google/uuid"
)

func GameToDomain(rg *rp.GameModel) (*domainModel.Game, error) {
	if rg == nil {
		return nil, errors.New("failed to mapping into repo model")
	}
	cells := rg.Board.ToTwoDimensionArray()
	domainBoard := domainModel.Board{
		Cells: cells,
	}
	turn := domainModel.Player{}
	winner := domainModel.Player{}
	firstPlayer := domainModel.Player{
		ID:   rg.FirstPlayerID,
		Icon: 1,
	}

	secondPlayer := domainModel.Player{
		ID:   rg.SecondPlayerID,
		Icon: 2,
	}

	switch rg.CurrentTurn {
	case "X":
		turn.Icon = 1
		turn.ID = rg.FirstPlayerID
	case "O":
		turn.Icon = 2
		turn.ID = rg.SecondPlayerID
	default:
		turn.Icon = 0
		turn.ID = ""
	}

	if rg.Winner != "" && rg.Winner == rg.FirstPlayerID {
		winner.Icon = 1
		winner.ID = rg.FirstPlayerID
	} else if rg.Winner != "" && rg.Winner == rg.SecondPlayerID {
		winner.Icon = 2
		winner.ID = rg.SecondPlayerID
	} else if rg.Winner == "draw" {
		winner.Icon = 0
		winner.ID = "draw"
	} else {
		winner.Icon = 0
		winner.ID = ""
	}

	return &domainModel.Game{
		ID:           rg.ID,
		Board:        domainBoard,
		FirstPlayer:  firstPlayer,
		SecondPlayer: secondPlayer,
		CurrentTurn:  turn,
		Status:       rg.Status,
		Winner:       winner,
	}, nil
}

func GameFromDomain(dm *domainModel.Game) (*rp.GameModel, error) {
	if dm == nil {
		return nil, errors.New("failed to mapping into repo model")
	}
	cells := dm.Board.ToOneDimensionArray()
	repoBoard := rp.BoardModel{
		Cells: cells,
	}

	turn := ""
	switch dm.CurrentTurn.Icon {
	case 1:
		turn = "X"
	case 2:
		turn = "O"
	default:
		turn = ""
	}

	winnerID := ""
	if dm.Status == domainModel.StatusGameOver && dm.Winner.ID != "" {
		winnerID = dm.Winner.ID
	} else if dm.Status == domainModel.StatusDraw {
		winnerID = "draw"
	}

	return &rp.GameModel{
		ID:             dm.ID,
		Board:          repoBoard,
		FirstPlayerID:  dm.FirstPlayer.ID,
		SecondPlayerID: dm.SecondPlayer.ID,
		CurrentTurn:    turn,
		Winner:         winnerID,
		Status:         dm.Status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func UserToDomain(rg *rp.UserModel) (*domainModel.User, error) {
	uu_id, err := uuid.Parse(rg.UUID)
	if err != nil {
		log.Println("failed to mapping user model")
	}
	return &domainModel.User{
		ID:       uu_id,
		Login:    rg.Login,
		Password: rg.Password,
	}, nil
}

func UserFromDomain(dm *domainModel.User) (*rp.UserModel, error) {
	uu_id := dm.ID.String()
	return &rp.UserModel{
		UUID:      uu_id,
		Login:     dm.Login,
		Password:  dm.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
