package service

import (
	"tictactoe/internal/domain"
)

type Move struct {
	Row   uint8
	Col   uint8
	Score int
}

type MinMaxServiceImpl struct {
}

func NewMinMaxService() *MinMaxServiceImpl {
	return &MinMaxServiceImpl{}
}
func checkWin(board domain.Board, player uint8) bool {
	for i := range 3 {
		if board.Cells[i][0] == player &&
			board.Cells[i][1] == player &&
			board.Cells[i][2] == player {
			return true
		}
	}

	for i := range 3 {
		if board.Cells[0][i] == player &&
			board.Cells[1][i] == player &&
			board.Cells[2][i] == player {
			return true
		}
	}

	if board.Cells[0][0] == player &&
		board.Cells[1][1] == player &&
		board.Cells[2][2] == player {
		return true
	}

	if board.Cells[0][2] == player &&
		board.Cells[1][1] == player &&
		board.Cells[2][0] == player {
		return true
	}

	return false
}

func (*MinMaxServiceImpl) GetBestMove(g domain.Game) (uint8, uint8) {
	bestMoves := minimax(g.Board, g.CurrentTurn.Icon, 10)
	return bestMoves.Row, bestMoves.Col
}

func minimax(board domain.Board, player uint8, depth int) Move {
	emptyCells := board.GetEmptyCells()
	if checkWin(board, 1) { // ai
		return Move{Score: 10 - depth}
	}
	if checkWin(board, 2) { // human
		return Move{Score: -10 + depth}
	}
	if len(emptyCells) == 0 {
		return Move{Score: 0}
	}

	moves := make([]Move, 0)

	for coords := range emptyCells {
		move := Move{
			Row: coords[0],
			Col: coords[1],
		}
		oldValue := board.Cells[coords[0]][coords[1]]
		board.Cells[coords[0]][coords[1]] = player
		if player == 2 { // ai
			result := minimax(board, 1, depth+1)
			move.Score = result.Score
		} else { // human
			result := minimax(board, 2, depth+1)
			move.Score = result.Score
		}
		board.Cells[coords[0]][coords[1]] = oldValue

		moves = append(moves, move)
	}

	bestMove := Move{}
	if player == 2 {
		bestScore := -10000
		for _, move := range moves {
			if move.Score > bestScore {
				bestScore = move.Score
				bestMove = move
			}
		}
	} else {
		bestScore := 10000
		for _, move := range moves {
			if move.Score < bestScore {
				bestScore = move.Score
				bestMove = move
			}
		}
	}

	return bestMove
}
