package web

import (
	"net/http"
	"tictactoe/internal/controller/handler"
	"tictactoe/internal/controller/middleware"
)

func NewRouter(handler handler.GameHandler, authHandler middleware.AuthHandler, userHandler handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/auth/signin", authHandler.SignIn)
	mux.HandleFunc("POST /api/auth/signup", authHandler.SignUp)
	mux.HandleFunc("POST /api/auth/refresh/access", authHandler.UpdateAccessToken)
	mux.HandleFunc("POST /api/auth/refresh/refresh", authHandler.UpdateRefreshToken)
	mux.HandleFunc("GET /api/", handler.MainPage)

	mux.HandleFunc("POST /api/game/ai", authHandler.Authenticate(handler.CreateGameWithAI))
	mux.HandleFunc("POST /api/game/multiplayer", authHandler.Authenticate(handler.CreateMultiplayerGame))
	mux.HandleFunc("POST /api/game/multiplayer/{join}", authHandler.Authenticate(handler.JoinToGame))
	mux.HandleFunc("GET /api/game/available", handler.GetWaitingGames)

	mux.HandleFunc("POST /api/game/{gameID}", authHandler.Authenticate(handler.PlayTurn))
	mux.HandleFunc("GET /api/game/{gameID}", authHandler.Authenticate(handler.GetGame))

	mux.HandleFunc("GET /api/user/me", authHandler.Authenticate(userHandler.GetUserInfoByAccessToken))
	mux.HandleFunc("GET /api/user/{userID}", userHandler.GetUserById)
	return mux
}
