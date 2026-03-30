package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/usecase/auth"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(a *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: a,
	}
}

func (au *AuthHandler) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Invalid login or password", http.StatusUnauthorized)
			return
		}
		userID, err := au.authService.SignIn(r.Context(), username, password)
		if err != nil {
			http.Error(w, "Unauthorization", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserById(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	return userID, ok
}

func (au *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	req := model.SignUpRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		au.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := au.authService.SignUp(r.Context(), req); err != nil {
		au.writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	au.writeJSON(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
		"code":    fmt.Sprint(http.StatusOK),
	})
}

func (au *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		au.writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	userID, err := au.authService.SignIn(r.Context(), username, password)
	if err != nil {
		log.Println(err)
		au.writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	au.writeJSON(w, http.StatusOK, map[string]string{
		"user_id": userID.String(),
		"code":    fmt.Sprint(http.StatusOK),
		"message": "Login successfully",
	})
}

func (au *AuthHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (au *AuthHandler) writeError(w http.ResponseWriter, status int, message string) {
	resp := response.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	}
	au.writeJSON(w, status, resp)
}
