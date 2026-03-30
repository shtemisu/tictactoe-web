package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/usecase/user"
)

type UserHandler struct {
	UserService *user.UserService
}

func NewUserHandler(us *user.UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if userID == "" {
		uh.writeError(w, http.StatusNotFound, "user not found")
	}
	userInfo, err := uh.UserService.GetUserInfo(r.Context(), userID)
	if err != nil {
		uh.writeError(w, http.StatusNotFound, err.Error())
	}
	uh.writeJSON(w, http.StatusOK, userInfo)
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
