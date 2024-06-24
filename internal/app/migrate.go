package app

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func init() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_URL")
	}
	fmt.Println("START TO MIGRAGE!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	//postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable

	databaseURL += "?sslmode=disable"

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		fmt.Printf("Db url: %s\n", databaseURL)
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		fmt.Printf("Migrate: pgdb is trying to connect, attempts left: %d\n", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: pgdb connect error: %s", err)
	}

	err = m.Up()
	defer func() { _, _ = m.Close() }()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("Migrate: no change")
		return
	}

	slog.Info("Migrate: up success")
}
