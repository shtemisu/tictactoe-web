package service

import (
	"context"
	"errors"
	"log"
	"tictactoe/internal/domain/model"
	db "tictactoe/internal/domain/repository"
	"tictactoe/internal/domain/service"
	datasource "tictactoe/internal/repository/dto"
)

type GameServiceImpl struct {
	repo   db.GameRepository
	minmax service.MinMaxService
}

// метод-конструктор.
func NewGameService(repo db.GameRepository, minmax service.MinMaxService) *GameServiceImpl {
	return &GameServiceImpl{
		repo:   repo,
		minmax: minmax,
	}
}

// возвращает gameID - сгенерированный uuid, записывает в db
func (gs *GameServiceImpl) CreateGameWithAI(ctx context.Context, playerID string) (string, error) {
	player := model.Player{
		ID:   playerID,
		Icon: 1,
	}
	game := model.InitGameWithAI(player)
	repoModel, err := datasource.GameFromDomain(game)
	if err != nil {
		return "", err
	}
	gs.repo.SaveGame(ctx, *repoModel)
	return game.ID, nil
}

func (gs *GameServiceImpl) CreateMultiplayerGame(ctx context.Context, firstPlayerID string) (string, error) {
	player := model.Player{
		ID:   firstPlayerID,
		Icon: 1,
	}
	game := model.InitMultiplayerGame(player)
	repoModel, err := datasource.GameFromDomain(game)
	if err != nil {
		return "", err
	}
	gs.repo.SaveGame(ctx, *repoModel)
	return game.ID, nil
}

func (gs *GameServiceImpl) JoinToGame(ctx context.Context, secondPlayerID string, gameID string) (*model.Game, error) {
	game, err := gs.repo.FindGameByWaitingStatus(ctx, gameID)
	if err != nil {
		return nil, err
	}
	domainModel, err := datasource.GameToDomain(game)
	if err != nil {
		return nil, err
	}
	domainModel.SecondPlayer.ID = secondPlayerID

	if err := gs.updateGameStatus(ctx, domainModel); err != nil {
		return nil, err
	}
	return domainModel, nil
}

func (gs *GameServiceImpl) GetTurn(ctx context.Context, gameID string) (uint8, error) {
	value, err1 := gs.repo.FindGameById(ctx, gameID)
	if err1 != nil {
		return 100, errors.New("Invalid game data")
	}
	domainModel, err2 := datasource.GameToDomain(value)
	if err2 != nil {
		return 100, err2
	}
	return domainModel.CurrentTurn.Icon, nil
}

func (gs *GameServiceImpl) GetGame(ctx context.Context, gameID string) (*model.Game, error) {
	value, err1 := gs.repo.FindGameById(ctx, gameID)

	if err1 != nil {
		return nil, errors.New("Invalid game data")
	}
	domainModel, err2 := datasource.GameToDomain(value)
	if err2 != nil {
		return nil, err2
	}
	return domainModel, nil
}
func (gs *GameServiceImpl) GetAllWaitingGames(ctx context.Context) ([]*model.Game, error) {
	games, err := gs.repo.FindAllWaitingGames(ctx)
	if err != nil {
		return nil, err
	}

	var domainGames []*model.Game
	for _, game := range games {
		domainGame, err := datasource.GameToDomain(game)
		if err != nil {
			continue
		}
		domainGames = append(domainGames, domainGame)
	}

	return domainGames, nil
}

func (gs *GameServiceImpl) updateGameStatus(ctx context.Context, game *model.Game) error {
	if game.SecondPlayer.ID != "" && game.FirstPlayer.ID != "" && game.Status == model.StatusWaiting {
		game.Status = model.StatusPlaying
		log.Println("Game status updated to playing, current turn remains with:", game.CurrentTurn.ID)
	}

	if game.CheckWinPlayer(game.CurrentTurn.Icon) {
		game.Winner = game.CurrentTurn
		game.Status = model.StatusGameOver
	} else if game.IsDraw() {
		game.Winner = model.Player{ID: "", Icon: 0}
		game.Status = model.StatusDraw
	}

	repoModel, err := datasource.GameFromDomain(game)
	if err != nil {
		return err
	}
	return gs.repo.UpdateGame(ctx, *repoModel)
}

// switchTurn меняет очередность хода
func (gs *GameServiceImpl) switchTurn(ctx context.Context, game *model.Game) error {
	if game.Status != model.StatusPlaying {
		return nil
	}

	if game.CurrentTurn.ID == game.FirstPlayer.ID {
		game.CurrentTurn = game.SecondPlayer
	} else {
		game.CurrentTurn = game.FirstPlayer
	}

	repoModel, err := datasource.GameFromDomain(game)
	if err != nil {
		return err
	}
	return gs.repo.UpdateGame(ctx, *repoModel)
}

// checkAndUpdateGameStatus проверяет завершение игры и обновляет статус
func (gs *GameServiceImpl) checkAndUpdateGameStatus(ctx context.Context, game *model.Game) error {
	gameOver := false

	if game.CheckWinPlayer(game.CurrentTurn.Icon) {
		game.Winner = game.CurrentTurn
		game.Status = model.StatusGameOver
		gameOver = true
	} else if game.IsDraw() {
		game.Winner = model.Player{ID: "", Icon: 0}
		game.Status = model.StatusDraw
		gameOver = true
	}

	if game.SecondPlayer.ID != "" && game.FirstPlayer.ID != "" && game.Status == model.StatusWaiting {
		game.Status = model.StatusPlaying
	}

	repoModel, err := datasource.GameFromDomain(game)
	if err != nil {
		return err
	}

	if err := gs.repo.UpdateGame(ctx, *repoModel); err != nil {
		return err
	}
	if gameOver {
		return nil
	}
	return nil
}

func (gs *GameServiceImpl) DoMove(ctx context.Context, gameID string, row uint8, col uint8) (bool, error) {
	currentGame, err1 := gs.repo.FindGameById(ctx, gameID)
	if err1 != nil {
		return false, errors.New("Invalid game data")
	}

	domainModel, err2 := datasource.GameToDomain(currentGame)
	if err2 != nil {
		return false, err2
	}

	if domainModel.Status != model.StatusPlaying {
		return false, errors.New("game is not active")
	}

	errMove := domainModel.MakeMove(row, col)
	if errMove != nil {
		return false, errMove
	}

	if err := gs.checkAndUpdateGameStatus(ctx, domainModel); err != nil {
		return false, err
	}

	if domainModel.Status == model.StatusPlaying {
		if err := gs.switchTurn(ctx, domainModel); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (gs *GameServiceImpl) GetNextMove(ctx context.Context, gameID string) (uint8, uint8, error) {
	game, err := gs.repo.FindGameById(ctx, gameID)
	if err != nil {
		return 100, 100, errors.New("game not found")
	}
	domainModel, err := datasource.GameToDomain(game)
	if err != nil {
		return 100, 100, err
	}
	x, y := gs.minmax.GetBestMove(*domainModel)
	return x, y, nil
}
