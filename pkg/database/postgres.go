package database

import (
	"fmt"
	"log/slog"
	"time"
	"ttages/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase(cfg config.Database, log *slog.Logger) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode,
	)

	var DB *sqlx.DB
	var err error

	for i := 1; i <= 5; i++ { // Попытки подключения
		DB, err = sqlx.Open("postgres", dsn)
		if err == nil {
			err = DB.Ping()
		}

		if err == nil {
			log.Info("Connected to database")
			return DB
		}

		log.Warn(fmt.Sprintf("try %d: cannot connect to db: %s", i, err))
		time.Sleep(5 * time.Second) // Ждём перед повторной попыткой
	}

	log.Error("Failed to connect to database after 5 tries")
	panic("Database connection failed")
}
