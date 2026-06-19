package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/controller/middleware"
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

// GetMyProfile returns the authenticated user's info: {id, login, games_played, wins}.
// The userID is provided by the auth middleware via the request context.
func (uh *UserHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserById(r.Context())
	if !ok {
		uh.writeError(w, http.StatusUnauthorized, "User unauthorized")
		return
	}
	userInfo, err := uh.UserService.GetUserInfo(r.Context(), userID.String())
	if err != nil {
		uh.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	uh.writeJSON(w, http.StatusOK, userInfo)
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
