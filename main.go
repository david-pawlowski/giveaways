package main

import (
	"github.com/david-pawlowski/giveaway/handlers"
	"github.com/david-pawlowski/giveaway/repository"
	"github.com/david-pawlowski/giveaway/service"
	"net/http"
)

func main() {
	store := repository.NewInMemoryStore()
	giveawayService := service.NewGiveawayService(store)
	giveawayHandler := handlers.NewGiveawayHandler(giveawayService)
	homeHandler := &handlers.HomeHandler{}

	mux := http.NewServeMux()
	mux.Handle("/", homeHandler)
	mux.Handle("/codes/", giveawayHandler)
	mux.Handle("/codes", giveawayHandler)

	http.ListenAndServe(":8080", mux)
}
