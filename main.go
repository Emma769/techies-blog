package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/emma769/techies-blog/internal/config"
	"github.com/emma769/techies-blog/internal/handler"
	"github.com/emma769/techies-blog/internal/middleware"
	"github.com/emma769/techies-blog/internal/repository/pgstore"
	"github.com/emma769/techies-blog/internal/services/article"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	config := config.New()

	store, err := pgstore.New(config)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	defer store.Close()

	logger.InfoContext(context.Background(), "db connected")

	mux := chi.NewRouter()

	mux.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Logger: logger,
	}))

	mux.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Logger: logger,
	}))

	mux.Use(middleware.EnableCORS(middleware.CORSConfig{
		AllowCredentials: false,
		Headers:          []string{"*"},
		Methods:          []string{"GET", "POST", "PUT", "DELETE"},
		Origins:          []string{"*"},
	}))

	article := article.NewService(store)

	services := &handler.Services{
		Article: article,
	}

	handler := handler.New(services)

	handler.Register(mux)

	port := config.LoadInt("PORT", 9000)

	logger.InfoContext(context.Background(), "server is starting", slog.Int("port", port))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
