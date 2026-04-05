package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/usecase/user"
	"tictactoe/pkg/jwt"
)

type UserHandler struct {
	UserService *user.UserService
	JwtProvider *jwt.JwtProvider
}

func NewUserHandler(us *user.UserService, jp *jwt.JwtProvider) *UserHandler {
	return &UserHandler{
		UserService: us,
		JwtProvider: jp,
	}
}

func (uh *UserHandler) GetUserInfoByAccessToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		uh.writeError(w, http.StatusBadRequest, "Bad request")
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := uh.JwtProvider.GetUUIDByToken(tokenStr, false)
	if err != nil {
		uh.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	userResp, err := uh.UserService.GetUserInfo(r.Context(), userID)
	if err != nil {
		uh.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	uh.writeJSON(w, http.StatusOK, userResp)
}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if userID == "" {
		uh.writeError(w, http.StatusBadRequest, "bad request")
		return
	}
	userInfo, err := uh.UserService.GetUserInfo(r.Context(), userID)
	if err != nil {
		uh.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	uh.writeJSON(w, http.StatusOK, userInfo)
}

func (uh *UserHandler) GetLeaderBoard(w http.ResponseWriter, r *http.Request) {
	limit := r.PathValue("limit")
	if limit == "" {
		uh.writeError(w, http.StatusBadRequest, "bad request")
		return
	}
	leaderBoard, err := uh.UserService.GetLeaderBoard(r.Context(), limit)
	if err != nil {
		uh.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	uh.writeJSON(w, http.StatusOK, leaderBoard)
}

func (uh *UserHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (uh *UserHandler) writeError(w http.ResponseWriter, status int, message string) {
	resp := response.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	}
	uh.writeJSON(w, status, resp)
}
