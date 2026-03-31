package user

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"log"
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/domain/repository/db"
	domainService "tictactoe/internal/domain/service"
	"tictactoe/internal/repository/db/dto"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo db.UserRepository
}

var _ domainService.UserService = (*UserService)(nil)

func NewUserService(userRepo db.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, req model.SignUpRequest) error {
	userExists, _ := u.userRepo.FindUserByLogin(ctx, req.Login)
	if userExists != nil {
		return errors.New("user with that login already exists")
	}

	hashPassword := sha256.Sum256([]byte(req.Password))
	user := &model.User{
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

func (u *UserService) GetUserBylogin(ctx context.Context, login string) (*model.User, error) {
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

func (u *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
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

func (s *UserService) GetUserInfo(ctx context.Context, ID string) (*response.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, ID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	gamesPlayed, wins, err := s.userRepo.GetUserStats(ctx, user.UUID)
	if err != nil {
		return nil, errors.New("failed to get user stats")
	}

	return &response.UserResponse{
		ID:          user.UUID,
		Login:       user.Login,
		GamesPlayed: gamesPlayed,
		Wins:        wins,
	}, nil
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
