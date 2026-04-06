package db

import (
	"context"
	"errors"
	"log"

	"tictactoe/internal/domain"
	domainRepo "tictactoe/internal/domain/repository"
	rp "tictactoe/internal/repository/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

var _ domainRepo.UserRepository = (*UserRepository)(nil)

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) FindUserByLogin(ctx context.Context, login string) (*rp.UserModel, error) {
	var u rp.UserModel
	err := r.pool.QueryRow(ctx, "SELECT id, login, password_hash FROM users WHERE login=$1", login).Scan(&u.UUID, &u.Login, &u.Password)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, ID string) (*rp.UserModel, error) {
	var u rp.UserModel
	err := r.pool.QueryRow(ctx, "SELECT id, login FROM users WHERE id=$1", ID).Scan(&u.UUID, &u.Login)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetUserStats(ctx context.Context, ID string) (gamesPlayed int, wins int, err error) {
	query := `
		SELECT
			COUNT(DISTINCT g.id) as games_played,
			COUNT(DISTINCT CASE WHEN g.winner = $1 THEN g.id END) as wins
		FROM games g
		WHERE g.firstPlayer_id = $1 OR g.secondPlayer_id = $1 
	`

	err = r.pool.QueryRow(ctx, query, ID).Scan(&gamesPlayed, &wins)
	if err != nil {
		return 0, 0, err
	}

	return gamesPlayed, wins, nil
}

func (r *UserRepository) GetLeaderBoard(ctx context.Context, limit string) ([]domain.LeaderBoardResponse, error) {
	var leaderBoard []domain.LeaderBoardResponse
	query := `
		SELECT u.id, ROUND(
			COUNT(CASE WHEN g.winner = u.id THEN 1 END)::DECIMAL /
			NULLIF(COUNT(DISTINCT g.id), 0) * 100,
			2
		) as winrate
		FROM users u
		LEFT JOIN games g ON 
			(g.firstPlayer_id = u.id OR g.secondPlayer_id = u.id)
			AND g.status IN ('game_over', 'draw')
		GROUP BY u.id
		HAVING COUNT(DISTINCT g.id) > 0
		ORDER BY winrate DESC
		LIMIT $1;
	`
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var leader domain.LeaderBoardResponse
		err := rows.Scan(&leader.UUID, &leader.Winrate)
		if err != nil {
			continue
		}
		leaderBoard = append(leaderBoard, leader)
	}
	return leaderBoard, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, u rp.UserModel) error {
	var uuid string
	err := r.pool.QueryRow(ctx,
		"INSERT INTO users(id, login, password_hash, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		u.UUID, u.Login, u.Password, u.CreatedAt, u.UpdatedAt).Scan(&uuid)
	if err != nil {
		log.Printf("%s\n", err)
		return errors.New("failed to save user in db")
	} else {
		log.Printf("user with ID: %s was save", uuid)
	}
	return nil
}

func (r *UserRepository) RemoveUser(ctx context.Context, user_id string) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.pool.Exec(ctx, query, user_id)
	if err != nil {
		return errors.New("failed to delete user")
	} else {
		log.Printf("user with ID: %s was delete", user_id)
	}
	if result.RowsAffected() == 0 {
		return errors.New("game not found")
	}
	return nil
}
