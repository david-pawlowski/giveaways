package main

import (
	"github.com/david-pawlowski/giveaway/handlers"
	"github.com/david-pawlowski/giveaway/repository"
	"github.com/david-pawlowski/giveaway/service"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	store := repository.InMemoryStore{}
	giveawayService := service.NewGiveawayService(&store)
	giveawayHandler := handlers.NewGiveawayHandler(giveawayService)

	mux := http.NewServeMux()
	mux.Handle("/", giveawayHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://main.d3ue8i9jys9zyg.amplifyapp.com/"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	http.ListenAndServe(":8080", corsHandler.Handler(mux))
}
