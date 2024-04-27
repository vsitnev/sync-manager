package app

import (
	log "log/slog"

	_"github.com/vsitnev/sync-manager/config"
	_"github.com/vsitnev/sync-manager/pkg/postgres"
)

func Run() {
	 log.Info("Starting application...")
	// cfg, err := config.NewConfig()
	// if err != nil {
	// 	log.Error("could not parse config: %w", err)
	// 	return
	// }

	// // init Db
	// db, err := postgres.New(cfg.DSN.Database)
	// if err != nil {
	// 	log.Error("app - Run - postgres.New: %w", err)
	// 	return
	// }

	// // init amqp
	// amqp, err := amqp.New(cfg.AMQP.URL)
	// if err != nil {
	// 	log.Error("app - Run - amqp.New: %w", err)
	// 	return
	// }


	
}
