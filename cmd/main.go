package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Souvik9205/go-api/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres dbname=ecom password=postgres sslmode=disable"),
		},
	}

	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("db connected", "dsn", cfg.db.dsn)

	api := applicattion{
		config: cfg,
		db:     conn,
	}

	h := api.mount()
	if err := api.run(h); err != nil {
		logger.Error("Server has failed to start, err: %s", err)
		os.Exit(1)
	}
}
