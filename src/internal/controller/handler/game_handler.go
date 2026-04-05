package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/controller/dto/mapper"
	"tictactoe/internal/controller/dto/request"
	"tictactoe/internal/controller/dto/response"
	"tictactoe/internal/controller/middleware"
	"tictactoe/internal/domain"
	"time"
)

type GameHandler struct {
	gameService domain.GameService
}

func NewGameHandler(gameService domain.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (h *GameHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *GameHandler) writeError(w http.ResponseWriter, status int, message string) {
	resp := response.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	}
	h.writeJSON(w, status, resp)
}

func (h *GameHandler) CreateGameWithAI(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserById(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "User unauthorized")
		return
	}
	gameID, err := h.gameService.CreateGameWithAI(r.Context(), userID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := response.CreateGameResponse{
		GameID:    gameID,
		Message:   "The game was created",
		CreatedAt: time.Now(),
	}
	h.writeJSON(w, http.StatusCreated, resp)
}

func (h *GameHandler) CreateMultiplayerGame(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserById(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "User unauthorized")
		return
	}
	gameID, err := h.gameService.CreateMultiplayerGame(r.Context(), userID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := response.CreateGameResponse{
		GameID:    gameID,
		Message:   "The game was created",
		CreatedAt: time.Now(),
	}
	h.writeJSON(w, http.StatusCreated, resp)
}

func (h *GameHandler) JoinToGame(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserById(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "User unauthorized")
		return
	}
	gameID := r.PathValue("join")
	if gameID == "" {
		h.writeError(w, http.StatusNotFound, "game_id is required")
		return
	}
	game, err := h.gameService.JoinToGame(r.Context(), userID.String(), gameID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	resp := mapper.DomainToResponse(game)
	h.writeJSON(w, http.StatusAccepted, resp)
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("gameID")
	if gameID == "" {
		h.writeError(w, http.StatusNotFound, "game_id is required")
		return
	}
	game, err := h.gameService.GetGame(r.Context(), gameID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "game not found")
		return
	}
	resp := mapper.DomainToResponse(game)
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *GameHandler) GetWaitingGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameService.GetAllWaitingGames(r.Context())
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	resp := mapper.DomainToWaitingGamesResponse(games)
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *GameHandler) GetGameHistoryByPlayerID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("playerID")
	gamesID, err := h.gameService.GetGameHistoryByPlayerID(r.Context(), userID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "games not found")
		return
	}
	h.writeJSON(w, http.StatusOK, gamesID)
}

func (h *GameHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Info []string `json:"info"`
	}{Info: []string{
		"Welcome! It is main page of tictactoe game",
		"Methods:",
		"GET /",
		"POST /game/{gameID}",
		"GET /game/{gameID}",
	}}
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *GameHandler) PlayTurn(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("gameID")
	if gameID == "" {
		http.Error(w, "game_id is required", http.StatusNotFound)
		return
	}

	moveReq := request.MakeMoveRequest{}

	if err := json.NewDecoder(r.Body).Decode(&moveReq); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if success, err := h.gameService.DoMove(r.Context(), gameID, uint8(moveReq.Row), uint8(moveReq.Col)); err != nil || !success {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	game, _ := h.gameService.GetGame(r.Context(), gameID)

	if game.SecondPlayer.ID == "minmax" && game.Status == domain.StatusPlaying {
		x, y, err := h.gameService.GetNextMove(r.Context(), gameID)
		if err != nil {
			h.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if success, err := h.gameService.DoMove(r.Context(), gameID, x, y); err != nil || !success {
			h.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		game, _ = h.gameService.GetGame(r.Context(), gameID)
	}

	resp := mapper.DomainToResponse(game)
	h.writeJSON(w, http.StatusAccepted, resp)
}
