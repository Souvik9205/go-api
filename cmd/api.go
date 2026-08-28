package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/Souvik9205/go-api/internal/adapters/postgresql/sqlc"
	"github.com/Souvik9205/go-api/internal/product"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type applicattion struct {
	config config
	db     *pgx.Conn
}

// mount
func (app *applicattion) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)              //help to build
	r.Use(middleware.ClientIPFromRemoteAddr) //tracing rate. limiting by real ip
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) //recover from crashs

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good!"))
	})

	productService := product.NewService(repo.New(app.db))
	productHandler := product.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)

	return r
}

// run
func (app *applicattion) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
