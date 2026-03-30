package web

import (
	"net/http"
	"tictactoe/internal/controller/handler"
	"tictactoe/internal/controller/middleware"
)

func NewRouter(handler handler.GameHandler, authHandler middleware.AuthHandler, userHandler handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signin", authHandler.SignIn)
	mux.HandleFunc("POST /signup", authHandler.SignUp)
	mux.HandleFunc("GET /", handler.MainPage)

	mux.HandleFunc("POST /game_ai", authHandler.Authenticate(handler.CreateGameWithAI))
	mux.HandleFunc("POST /multiplayer", authHandler.Authenticate(handler.CreateMultiplayerGame))
	mux.HandleFunc("POST /multiplayer/{join}", authHandler.Authenticate(handler.JoinToGame))
	mux.HandleFunc("GET /available", handler.GetWaitingGames)

	mux.HandleFunc("POST /game/{gameID}", authHandler.Authenticate(handler.PlayTurn))
	mux.HandleFunc("GET /game/{gameID}", authHandler.Authenticate(handler.GetGame))

	mux.HandleFunc("GET /user/{userID}", userHandler.GetUserById)
	return mux
}
