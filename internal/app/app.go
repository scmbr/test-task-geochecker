package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/config"
	"github.com/scmbr/test-task-geochecker/internal/delivery/http/handler"
	"github.com/scmbr/test-task-geochecker/internal/repository"
	"github.com/scmbr/test-task-geochecker/internal/server"
	"github.com/scmbr/test-task-geochecker/internal/service"
	"github.com/scmbr/test-task-geochecker/pkg/cache"
	"github.com/scmbr/test-task-geochecker/pkg/database/postgres"
	"github.com/scmbr/test-task-geochecker/pkg/logger"
	"github.com/scmbr/test-task-geochecker/pkg/redis"
)

func Run(configsDir string) {
	logger.Init()
	cfg, err := config.Init(configsDir)
	if err != nil {
		logger.Error("config initialization error", err)
		os.Exit(1)
	}
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.Name,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		logger.Error("database initialization error", err)
		os.Exit(1)
	}
	logger.Info("database connected successfully")
	redisClient, err := redis.NewRedis(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
	})
	if err != nil {
		logger.Error("failed to initialize redis:%s", err)
	}
	cacheProvider := cache.NewRedisCache(redisClient)
	repository := repository.NewRepository(db)
	service := service.NewService(service.Deps{
		Repos:        repository,
		RadiusMeters: cfg.SearchRadius,
		ApiKeySecret: cfg.ApiKeySecret,
		Cache:        cacheProvider,
	})
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get generic database interface: %v", err)
		os.Exit(1)
	}
	handler := handler.NewHandler(service, sqlDB, redisClient)
	server := server.NewServer(cfg, handler.Init())
	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error occurred while running server", err)
		}
	}()

	logger.Info("server started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logger.Error("failed to stop server", err)
	}

}
