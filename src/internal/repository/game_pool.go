package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	db "tictactoe/internal/domain/repository"
	rp "tictactoe/internal/repository/model"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ db.GameRepository = (*GameRepositoryImpl)(nil)

type GameRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) *GameRepositoryImpl {
	return &GameRepositoryImpl{pool: pool}
}

func (r *GameRepositoryImpl) FindGameById(ctx context.Context, id string) (*rp.GameModel, error) {
	var g rp.GameModel

	query := `SELECT id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at 
              FROM games WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game not found")
		}
		log.Printf("Error finding game: %v", err)
		return nil, err
	}

	return &g, nil
}

func (r *GameRepositoryImpl) FindGameByWaitingStatus(ctx context.Context, id string) (*rp.GameModel, error) {
	var g rp.GameModel

	query := `SELECT id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at 
              FROM games WHERE status = 'waiting' AND id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
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
		return nil, errors.New("failed to scan games")
	}
	return games, nil
}

func (r *GameRepositoryImpl) GetGameHistoryByPlayerID(ctx context.Context, playerID string) ([]string, error) {
	var gamesID []string
	query := "SELECT id FROM games WHERE (status != 'waiting' AND status != 'playing') AND ($1 IN (firstPlayer_id, secondPlayer_id))"
	rows, err := r.pool.Query(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var gameID string
		err := rows.Scan(&gameID)
		if err != nil {
			continue
		}
		gamesID = append(gamesID, gameID)
	}
	return gamesID, nil
}

func (r *GameRepositoryImpl) SaveGame(ctx context.Context, g rp.GameModel) error {
	var id string

	query := `INSERT INTO games(id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at) 
              VALUES($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id`

	err := r.pool.QueryRow(ctx, query,
		g.ID,
		g.Board.Cells,
		g.FirstPlayerID,
		g.SecondPlayerID, // sql.NullString
		g.CurrentTurn,    // sql.NullString
		g.Status,
		g.Winner, // sql.NullString
	).Scan(&id)

	if err != nil {
		log.Printf("Error saving game: %v", err)
		return fmt.Errorf("failed to save game in db: %w", err)
	}

	log.Printf("game with ID: %s was saved", id)
	return nil
}

func (r *GameRepositoryImpl) UpdateGame(ctx context.Context, g rp.GameModel) error {
	query := `UPDATE games SET 
                board = $1, 
                firstPlayer_id = $2, 
                secondPlayer_id = $3, 
                winner = $4, 
                status = $5, 
                current_turn = $6, 
                updated_at = NOW() 
              WHERE id = $7`

	result, err := r.pool.Exec(ctx, query,
		g.Board.Cells,
		g.FirstPlayerID,
		g.SecondPlayerID,
		g.Winner,
		g.Status,
		g.CurrentTurn,
		g.ID,
	)

	if err != nil {
		log.Printf("Error updating game: %v", err)
		return fmt.Errorf("failed to update game: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("game with id %s not found", g.ID)
	}

	log.Printf("game with ID: %s was updated", g.ID)
	return nil
}
