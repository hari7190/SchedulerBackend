package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"hari.foo/schedulerbackend/internal/database"
	"hari.foo/schedulerbackend/internal/event"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("Initiating Scheduler API")

	// DB connection
	db, err := database.New(database.Config{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "pass",
		DBName:   "scheduler",
		SSLMode:  "disable",
	})

	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	mux := http.NewServeMux()

	eventStore := event.NewMySQLStore(db)
	eventService := event.NewService(eventStore)
	eventHandler := event.NewHandler(eventService)
	eventHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	srv.ListenAndServe()
	logger.Info("starting server", "addr", srv.Addr)

}

type DummyStore struct{}

func (d *DummyStore) Create(e event.Event) error             { return nil }
func (d *DummyStore) GetById(id string) (event.Event, error) { return event.Event{}, nil }
