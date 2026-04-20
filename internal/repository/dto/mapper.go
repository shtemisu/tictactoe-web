package dto

import (
	"database/sql"
	"errors"
	"log"
	"tictactoe/internal/domain"
	rp "tictactoe/internal/repository/model"
	"time"

	"github.com/google/uuid"
)

func GameToDomain(rg *rp.GameModel) (*domain.Game, error) {
	if rg == nil {
		return nil, errors.New("failed to mapping into repo model")
	}

	cells := rg.Board.ToTwoDimensionArray()
	domainBoard := domain.Board{
		Cells: cells,
	}

	turn := domain.Player{}
	winner := domain.Player{}

	firstPlayer := domain.Player{
		ID:   rg.FirstPlayerID,
		Icon: 1,
	}

	secondPlayer := domain.Player{
		ID:   rg.SecondPlayerID.String, // может быть пустой строкой
		Icon: 2,
	}

	// Обработка CurrentTurn (может быть NULL)
	if rg.CurrentTurn.Valid {
		if rg.CurrentTurn.String == rg.FirstPlayerID {
			turn.Icon = 1
			turn.ID = rg.FirstPlayerID
		} else if rg.CurrentTurn.String == rg.SecondPlayerID.String {
			turn.Icon = 2
			turn.ID = rg.SecondPlayerID.String
		}
	} else {
		// NULL - ход не определен
		turn.Icon = 0
		turn.ID = ""
	}

	// Обработка Winner (может быть NULL!)
	if rg.Winner.Valid {
		winnerID := rg.Winner.String

		if winnerID == rg.FirstPlayerID {
			winner.Icon = 1
			winner.ID = rg.FirstPlayerID
		} else if winnerID == rg.SecondPlayerID.String {
			winner.Icon = 2
			winner.ID = rg.SecondPlayerID.String
		} else if winnerID == "draw" {
			winner.Icon = 0
			winner.ID = "draw"
		} else {
			winner.Icon = 0
			winner.ID = ""
		}
	} else {
		// NULL - победителя нет
		winner.Icon = 0
		winner.ID = ""
	}

	return &domain.Game{
		ID:           rg.ID,
		Board:        domainBoard,
		FirstPlayer:  firstPlayer,
		SecondPlayer: secondPlayer,
		CurrentTurn:  turn,
		Status:       rg.Status,
		Winner:       winner,
	}, nil
}

func GameFromDomain(dm *domain.Game) (*rp.GameModel, error) {
	if dm == nil {
		return nil, errors.New("failed to mapping into repo model")
	}

	cells := dm.Board.ToOneDimensionArray()
	repoBoard := rp.BoardModel{
		Cells: cells,
	}

	// CurrentTurn - может быть NULL
	currentTurn := sql.NullString{Valid: false}
	if dm.CurrentTurn.ID != "" && dm.CurrentTurn.ID != "00000000-0000-0000-0000-000000000000" {
		currentTurn = sql.NullString{
			String: dm.CurrentTurn.ID,
			Valid:  true,
		}
	}

	// Winner - может быть NULL
	winner := sql.NullString{Valid: false}
	if dm.Winner.ID != "" && dm.Winner.ID != "00000000-0000-0000-0000-000000000000" && dm.Winner.ID != "draw" {
		winner = sql.NullString{
			String: dm.Winner.ID,
			Valid:  true,
		}
	} else if dm.Winner.ID == "draw" {
		winner = sql.NullString{
			String: "draw",
			Valid:  true,
		}
	}

	// SecondPlayer - может быть NULL
	secondPlayerID := sql.NullString{Valid: false}
	if dm.SecondPlayer.ID != "" && dm.SecondPlayer.ID != "00000000-0000-0000-0000-000000000000" {
		secondPlayerID = sql.NullString{
			String: dm.SecondPlayer.ID,
			Valid:  true,
		}
	}

	return &rp.GameModel{
		ID:             dm.ID,
		Board:          repoBoard,
		FirstPlayerID:  dm.FirstPlayer.ID,
		SecondPlayerID: secondPlayerID,
		CurrentTurn:    currentTurn,
		Winner:         winner,
		Status:         dm.Status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func UserToDomain(rg *rp.UserModel) (*domain.User, error) {
	uu_id, err := uuid.Parse(rg.UUID)
	if err != nil {
		log.Println("failed to mapping user model")
	}
	return &domain.User{
		ID:       uu_id,
		Login:    rg.Login,
		Password: rg.Password,
	}, nil
}

func UserFromDomain(dm *domain.User) (*rp.UserModel, error) {
	uu_id := dm.ID.String()
	return &rp.UserModel{
		UUID:      uu_id,
		Login:     dm.Login,
		Password:  dm.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
