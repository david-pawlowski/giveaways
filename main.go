package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	godotenv.Load(".env")
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

func setupServer(cfg *Config, handler http.Handler) *http.Server {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      corsMiddleware.Handler(handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config.")
	}

	store := repository.NewInMemoryStore()
	givServ, err := service.NewGiveawayService(store)
	if err != nil {
		log.Fatalf("Error connection to database failed: %v", err)
	}
	givHan := handlers.NewGiveawayHandler(givServ)

	mux := http.NewServeMux()
	mux.Handle("/", givHan)

	server := setupServer(cfg, mux)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Couldn't shutdown server: %v", err)
		}
		close(done)
	}()

	log.Printf("Server is starting on port %s...", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on port %s: %v", cfg.Port, err)
	}
	<-done
	log.Println("Server stopped")
}
