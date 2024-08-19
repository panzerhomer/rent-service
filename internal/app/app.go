package app

import (
	"avito/internal/config"
	"avito/internal/handlers"
	"avito/internal/repository"
	"avito/internal/server"
	"avito/internal/services"
	manager "avito/pkg/jwt"
	logging "avito/pkg/logger"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.Config) {
	log := logging.NewLogger(cfg.Env)

	log.Info(
		"starting avito-backend",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.PostgresDB)
	//
	connString = "postgres://postgres:admin@localhost:5432/go-avito-db?sslmode=disable"
	//
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Error("failed to init storage: " + err.Error())
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	log.Info(
		"connecting to postgres",
		slog.Int("port: ", cfg.Database.Port),
	)
	//
	cfg.Secret = "fdf7@@Ey932h7!9f"
	//
	tokenManager, err := manager.NewManager(cfg.Secret)
	if err != nil {
		log.Error("failed to init tokenManager: " + err.Error())
		return
	}

	userRepo := repository.NewUserRepository(pool, log)
	userService := services.NewUserService(userRepo, tokenManager, log)
	userHandler := handlers.NewUserHandler(userService, tokenManager, log)

	// flatRepo := repository.NewFlatRepository(pool, log)
	// flatService := services.NewFlatService(flatRepo, log)
	// flatHandler := handlers.New

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/dummyLogin", userHandler.DummyLogin)
	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)

	//
	cfg.App.Port = 8000
	//
	srv := server.NewServer(cfg, r)

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("error occurred while running http server: ", err.Error())
		}
	}()

	log.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Error("failed to stop server: %v", err)
	}
}
