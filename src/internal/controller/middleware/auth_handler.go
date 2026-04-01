package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/usecase/auth"
	"tictactoe/pkg/jwt"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *auth.AuthService
	jwtProvider *jwt.JwtProvider
}

func NewAuthHandler(a *auth.AuthService, jp *jwt.JwtProvider) *AuthHandler {
	return &AuthHandler{
		authService: a,
		jwtProvider: jp,
	}
}

func (au *AuthHandler) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			au.writeError(w, http.StatusBadRequest, "Bad request")
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		if !au.jwtProvider.ValidateAccessToken(tokenStr) {
			au.writeError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		userID, err := au.jwtProvider.GetUUIDByToken(tokenStr, false)
		if err != nil {
			au.writeError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			au.writeError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userUUID)
		next(w, r.WithContext(ctx))
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
	req := jwt.JwtRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		au.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	jwtResponse, err := au.authService.SignIn(r.Context(), req)
	if err != nil {
		au.writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	au.writeJSON(w, http.StatusOK, jwtResponse)
}
func (au *AuthHandler) UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
	req := jwt.RefreshJwtRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		au.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	jwtResponse, err := au.authService.UpdateAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		au.writeError(w, http.StatusUnauthorized, err.Error())
	}
	au.writeJSON(w, http.StatusOK, jwtResponse)
}

func (au *AuthHandler) UpdateRefreshToken(w http.ResponseWriter, r *http.Request) {
	req := jwt.RefreshJwtRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		au.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	jwtResponse, err := au.authService.UpdateRefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		au.writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	au.writeJSON(w, http.StatusOK, jwtResponse)
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
