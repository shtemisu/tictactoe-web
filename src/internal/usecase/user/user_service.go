package user

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"log"
	"tictactoe/internal/domain"
	db "tictactoe/internal/domain/repository"
	"tictactoe/internal/repository/dto"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo db.UserRepository
}

var _ domain.UserService = (*UserService)(nil)

func NewUserService(userRepo db.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, req domain.SignUpRequest) error {
	userExists, _ := u.userRepo.FindUserByLogin(ctx, req.Login)
	if userExists != nil {
		return errors.New("user with that login already exists")
	}

	hashPassword := sha256.Sum256([]byte(req.Password))
	user := &domain.User{
		ID:       uuid.New(),
		Login:    req.Login,
		Password: hex.EncodeToString(hashPassword[:]),
	}
	userRepoModel, err := dto.UserFromDomain(user)
	if err != nil {
		return err
	}
	return u.userRepo.SaveUser(ctx, *userRepoModel)
}

func (u *UserService) GetUserBylogin(ctx context.Context, login string) (*domain.User, error) {
	user, err := u.userRepo.FindUserByLogin(ctx, login)
	if err != nil {
		return nil, errors.New("user not found or not exists")
	}

	userDomain, err1 := dto.UserToDomain(user)
	if err1 != nil {
		return nil, errors.New("failed to mapping user model")
	}
	return userDomain, nil
}

func (u *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepo.FindUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found or not exists")
	}
	userDomain, err1 := dto.UserToDomain(user)
	if err1 != nil {
		return nil, errors.New("failed to mapping user model")
	}
	return userDomain, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, ID string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, ID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	gamesPlayed, wins, err := s.userRepo.GetUserStats(ctx, user.UUID)
	if err != nil {
		return nil, errors.New("failed to get user stats")
	}

	return &domain.UserResponse{
		ID:          user.UUID,
		Login:       user.Login,
		GamesPlayed: gamesPlayed,
		Wins:        wins,
	}, nil
}

func (s *UserService) GetLeaderBoard(ctx context.Context, limit string) ([]domain.LeaderBoardResponse, error) {
	leaderBoard, err := s.userRepo.GetLeaderBoard(ctx, limit)
	if err != nil {
		return nil, err
	}
	return leaderBoard, nil
}

func (u *UserService) ValidateCredentials(ctx context.Context, login string, password string) (uuid.UUID, error) {
	user, err := u.GetUserBylogin(ctx, login)
	if err != nil {
		log.Println("In user_service: ", err)
		return uuid.Nil, err
	}
	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])
	if subtle.ConstantTimeCompare([]byte(passwordHash), []byte(user.Password)) != 1 {
		return uuid.Nil, errors.New("invalid login or password")
	}
	return user.ID, nil
}
