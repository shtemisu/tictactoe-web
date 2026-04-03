package jwt

import (
	"tictactoe/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type JwtProvider struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewJwtProvider(accessSecret, refreshSecret string) *JwtProvider {
	return &JwtProvider{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (jp *JwtProvider) GenerateAccessToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"uuid": user.ID,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jp.accessSecret)
}

func (jp *JwtProvider) GenerateRefreshToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"uuid": user.ID,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jp.refreshSecret)
}

func (jp *JwtProvider) ValidateAccessToken(token string) bool {
	return jp.validate(token, jp.accessSecret)
}

func (jp *JwtProvider) ValidateRefreshToken(token string) bool {
	return jp.validate(token, jp.refreshSecret)
}

func (jp *JwtProvider) validate(tokenStr string, secret []byte) bool {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	return err == nil && token.Valid
}

func (jp *JwtProvider) GetUUIDByToken(tokenStr string, isRefresh bool) (string, error) {
	secret := jp.accessSecret
	if isRefresh {
		secret = jp.refreshSecret
	}
	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["uuid"].(string), nil
	}
	return "", jwt.ErrTokenInvalidClaims
}
