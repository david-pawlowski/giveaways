package main

import (
	"github.com/david-pawlowski/giveaway/handlers"
	"github.com/david-pawlowski/giveaway/repository"
	"github.com/david-pawlowski/giveaway/service"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func initDotEnv() {
	// TODO: dotenv should be local only
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func main() {
	initDotEnv()
	store := repository.InMemoryStore{}
	givServ := service.NewGiveawayService(&store)
	givHan := handlers.NewGiveawayHandler(givServ)

	mux := http.NewServeMux()
	mux.Handle("/", givHan)

	furl := os.Getenv("frontend_url")
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{furl},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	http.ListenAndServe(":8080", corsHandler.Handler(mux))
}
