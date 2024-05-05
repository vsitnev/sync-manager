package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/vsitnev/sync-manager/config"
	"github.com/vsitnev/sync-manager/internal/repository"
	"github.com/vsitnev/sync-manager/internal/service"
	v1 "github.com/vsitnev/sync-manager/internal/transport/http/v1"
	"github.com/vsitnev/sync-manager/pkg/amqp"
	"github.com/vsitnev/sync-manager/pkg/httpserver"
	"github.com/vsitnev/sync-manager/pkg/postgres"
)

// @title           Sync Manager Service
// @version         1.0
// @description     This is a service for sync services.

// @contact.name   Vladislav Sitnev
// @contact.email  vsitnev@yandex.ru

// @host      localhost:8080
// @BasePath  /
func Run() {
	slog.Info("Starting application...")
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("could not parse config: %w", err)
		return
	}

	// init db
	db, err := postgres.New(cfg.DSN.Database)
	if err != nil {
		slog.Error("app - Run - postgres.New: %w", err)
		return
	}

	// init amqp
	amqp, err := amqp.New(cfg.DSN.Amqp)
	if err != nil {
		slog.Error("app - Run - amqp.New: %w", err)
		return
	}

	// init repository
	reps := repository.NewRepositories(db)

	// init services
	services := service.NewServices(service.ServiceDeps{
		Reps: reps,
		Amqp: amqp,
	})

	// gin handler
	handler := gin.New()
	v1.NewRouter(handler, services)

	// http server
	slog.Info("Starting http server...")
	slog.Debug("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// waiting signal
	slog.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("app - Run - signal.Notify: %s", s.String())
	case err := <-httpServer.Notify():
		slog.Error("app - Run - httpServer.Notify: %w", err)
	}

	// graceful shutdown
	slog.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		slog.Error("app - Run - httpServer.Shutdown: %w", err)
	}
}
