package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/david-pawlowski/giveaway/handlers"
	"github.com/david-pawlowski/giveaway/repository"
	"github.com/david-pawlowski/giveaway/service"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type Config struct {
	Port        string
	FrontendURL string
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("Error loading .env: %w", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		DEFAULT_PORT := "8080"
		port = DEFAULT_PORT
	}
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		DEFAULT_FRONTEND_URL := "localhost:3000"
		frontendURL = DEFAULT_FRONTEND_URL
	}
	return &Config{
		Port:        port,
		FrontendURL: frontendURL,
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config.")
	}

	store := repository.NewInMemoryStore()
	givServ, err := service.NewGiveawayService(store)
	if err != nil {
		log.Fatal("Error connection to database failed.")
	}
	givHan := handlers.NewGiveawayHandler(givServ)

	mux := http.NewServeMux()
	mux.Handle("/", givHan)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	http.ListenAndServe(":"+cfg.Port, corsHandler.Handler(mux))
}
