package di

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"tictactoe/config"
	web "tictactoe/internal/controller"
	"tictactoe/internal/controller/handler"
	"tictactoe/internal/controller/middleware"
	domainRepo "tictactoe/internal/domain/repository"
	domainService "tictactoe/internal/domain/service"
	db "tictactoe/internal/repository/db"
	"tictactoe/internal/usecase/auth"
	usecases "tictactoe/internal/usecase/service"
	"tictactoe/internal/usecase/user"
	"tictactoe/pkg/jwt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func Module() fx.Option {
	module := fx.Module("tictactoe",
		fx.Provide(
			config.NewConfig,
			NewPostgresPool,
			db.NewGameRepository,
			fx.Annotate(
				db.NewGameRepository,
				fx.As(new(domainRepo.GameRepository)),
			),
			db.NewUserRepository,
			fx.Annotate(
				db.NewUserRepository,
				fx.As(new(domainRepo.UserRepository)),
			),
			usecases.NewMinMaxService,
			fx.Annotate(
				usecases.NewMinMaxService,
				fx.As(new(domainService.MinMaxService)),
			),
			usecases.NewGameService,
			fx.Annotate(
				usecases.NewGameService,
				fx.As(new(domainService.GameService)),
			),
			user.NewUserService,
			fx.Annotate(
				user.NewUserService,
				fx.As(new(domainService.UserService)),
			),
			fx.Annotate(func(cfg *config.Config) *jwt.JwtProvider {
				return jwt.NewJwtProvider(cfg.JWT_Access_Secret, cfg.JWT_Refresh_Secret)
			}),
			auth.NewAuthService,
			middleware.NewAuthHandler,
			handler.NewUserHandler,
			handler.NewGameHandler,
			NewHTTPServer,
		),
		fx.Invoke(InitDatabase),
		fx.Invoke(func(*http.Server) {}),
	)
	return module
}

func NewHTTPServer(lc fx.Lifecycle, handler *handler.GameHandler, authHandler *middleware.AuthHandler, userHandler *handler.UserHandler) *http.Server {
	router := web.NewRouter(*handler, *authHandler, *userHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Println("starting HTTP server on", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("stopping HTTP server on", srv.Addr)
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func NewPostgresPool(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)

	log.Printf("Connecting to database: %s", dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Printf("Database ping failed: %v", err)
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})
	return pool, nil
}

func InitDatabase(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), initGameTable)
	if err != nil {
		log.Printf("failed to create table with games in db %v", err)
		return err
	}
	_, err1 := pool.Exec(context.Background(), initUserTable)
	if err1 != nil {
		log.Printf("failed to create table with users in db: %v", err)
		return err
	}
	log.Println("Таблицы базы данных проверены/созданы")
	return nil
}
