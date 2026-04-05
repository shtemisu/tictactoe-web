package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

const (
	StatusPlaying  string = "playing"
	StatusDraw     string = "draw"
	StatusGameOver string = "game_over"
	StatusWaiting  string = "waiting"
)

type Player struct {
	ID   string
	Icon uint8
}

type Game struct {
	ID           string
	Board        Board
	FirstPlayer  Player
	SecondPlayer Player
	CurrentTurn  Player
	Status       string
	Winner       Player
}

type GameService interface {
	CreateGameWithAI(ctx context.Context, playerID string) (string, error)
	CreateMultiplayerGame(ctx context.Context, firstPlayerID string) (string, error)
	JoinToGame(ctx context.Context, secondPlayerID string, gameID string) (*Game, error)
	GetGame(ctx context.Context, gameID string) (*Game, error)
	GetAllWaitingGames(ctx context.Context) ([]*Game, error)
	GetGameHistoryByPlayerID(ctx context.Context, playerID string) (*[]string, error)
	GetTurn(ctx context.Context, gameID string) (uint8, error)
	DoMove(ctx context.Context, gameID string, row uint8, col uint8) (bool, error)
	GetNextMove(ctx context.Context, gameID string) (uint8, uint8, error)
}

func InitGameWithAI(firstPlayer Player) *Game {
	minmax := Player{
		ID:   "minmax",
		Icon: 2, // O
	}
	return &Game{
		ID:           uuid.New().String(),
		Board:        InitBoard(),
		FirstPlayer:  firstPlayer,
		SecondPlayer: minmax,
		CurrentTurn:  firstPlayer,
		Status:       StatusPlaying,
		Winner:       Player{},
	}
}

func InitMultiplayerGame(firstPlayer Player) *Game {
	return &Game{
		ID:           uuid.New().String(),
		Board:        InitBoard(),
		FirstPlayer:  firstPlayer,
		SecondPlayer: Player{},
		CurrentTurn:  firstPlayer,
		Status:       StatusWaiting,
		Winner:       Player{},
	}
}

func (g *Game) IsDraw() bool {
	for i := range 3 {
		for j := range 3 {
			if g.Board.Cells[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func (g *Game) CheckWinPlayer(player uint8) bool {
	for i := range 3 {
		if g.Board.Cells[i][0] == player &&
			g.Board.Cells[i][1] == player &&
			g.Board.Cells[i][2] == player {
			return true
		}
	}
	for i := range 3 {
		if g.Board.Cells[0][i] == player &&
			g.Board.Cells[1][i] == player &&
			g.Board.Cells[2][i] == player {
			return true
		}
	}
	if g.Board.Cells[0][0] == player &&
		g.Board.Cells[1][1] == player &&
		g.Board.Cells[2][2] == player {
		return true
	}

	if g.Board.Cells[0][2] == player &&
		g.Board.Cells[1][1] == player &&
		g.Board.Cells[2][0] == player {
		return true
	}

	return false
}

func (g *Game) CanMakeMove(row, col uint8) error {
	if g.Status != StatusPlaying {
		return errors.New("Game is over")
	} else if g.Board.Cells[row][col] != 0 {
		return errors.New("Cell is occupied")
	} else if row > 2 || col > 2 {
		return errors.New("Invalid value for row or col")
	}
	return nil
}

func (g *Game) MakeMove(row, col uint8) error {
	if err := g.CanMakeMove(row, col); err != nil {
		return err
	}

	g.Board.Cells[row][col] = g.CurrentTurn.Icon
	return nil
}
