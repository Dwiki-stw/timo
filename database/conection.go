package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"timo/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnection(conf config.DB) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name, conf.SslMode,
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	cfg.MaxConns = 100
	cfg.MinConns = 10
	cfg.MaxConnIdleTime = 3 * time.Minute
	cfg.MaxConnLifetime = 60 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	fmt.Println("database connected")

	return pool
}
