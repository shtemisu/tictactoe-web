package db

import (
	"errors"
	"fmt"
	"log"
	rp "tictactoe/internal/repository/model"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GameRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) *GameRepositoryImpl {
	return &GameRepositoryImpl{pool: pool}
}

func (r *GameRepositoryImpl) FindGameById(ctx context.Context, id string) (*rp.GameModel, error) {
	var g rp.GameModel
	err := r.pool.QueryRow(ctx, "SELECT * FROM games WHERE id=$1", id).Scan(&g.ID, &g.Board.Cells, &g.FirstPlayerID, &g.SecondPlayerID,
		&g.CurrentTurn, &g.Status, &g.Winner, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GameRepositoryImpl) FindGameByWaitingStatus(ctx context.Context, id string) (*rp.GameModel, error) {
	var g rp.GameModel
	err := r.pool.QueryRow(ctx, "SELECT * FROM games WHERE status='waiting' AND id=$1", id).Scan(&g.ID, &g.Board.Cells,
		&g.FirstPlayerID, &g.SecondPlayerID, &g.CurrentTurn, &g.Status, &g.Winner, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GameRepositoryImpl) FindAllWaitingGames(ctx context.Context) ([]*rp.GameModel, error) {
	var games []*rp.GameModel

	query := `SELECT id, board, firstplayer_id, secondplayer_id, current_turn, status, winner, created_at, updated_at 
	          FROM games WHERE status='waiting' ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g rp.GameModel
		err := rows.Scan(
			&g.ID,
			&g.Board.Cells,
			&g.FirstPlayerID,
			&g.SecondPlayerID,
			&g.CurrentTurn,
			&g.Status,
			&g.Winner,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {

			continue
		}
		games = append(games, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("failed to iterate games")
	}
	return games, nil
}

func (r *GameRepositoryImpl) SaveGame(ctx context.Context, g rp.GameModel) error {
	var id string
	fmt.Println(g.SecondPlayerID)
	query := "INSERT INTO games(id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id"
	err := r.pool.QueryRow(ctx, query, g.ID, g.Board.Cells, g.FirstPlayerID, g.SecondPlayerID, g.CurrentTurn, g.Status, g.Winner).Scan(&id)
	if err != nil {
		log.Printf("%s\n", err)
		return errors.New("failed to save game in db")
	} else {
		log.Printf("game with ID: %s was save", id)
	}
	return nil
}

func (r *GameRepositoryImpl) UpdateGame(ctx context.Context, g rp.GameModel) error {
	query := "UPDATE games SET board = $1, firstPlayer_id = $2, secondPlayer_id = $3, winner = $4, status = $5, current_turn = $6, updated_at = NOW() WHERE id = $7"
	result, err := r.pool.Exec(ctx, query, g.Board.Cells, g.FirstPlayerID, g.SecondPlayerID, g.Winner, g.Status, g.CurrentTurn, g.ID)
	if err != nil {
		return errors.New("failed to update game")
	} else {
		log.Printf("game with ID: %s was update", g.ID)
	}
	if result.RowsAffected() == 0 {
		return errors.New("game not found")
	}
	return nil
}

func (r *GameRepositoryImpl) RemoveGame(ctx context.Context, gameID string) error {
	query := "DELETE FROM games WHERE id = $1"
	result, err := r.pool.Exec(ctx, query, gameID)
	if err != nil {
		return errors.New("failed to delete game")
	} else {
		log.Printf("game with ID: %s was delete", gameID)
	}
	if result.RowsAffected() == 0 {
		return errors.New("game not found")
	}
	return nil
}
