package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/david-pawlowski/giveaway/models"
	"github.com/david-pawlowski/giveaway/repository"
	"github.com/david-pawlowski/giveaway/service"
)

type GiveawayHandler struct {
	giveaway service.GiveawayService
}

func NewGiveawayHandler(g service.GiveawayService) *GiveawayHandler {
	return &GiveawayHandler{
		giveaway: g,
	}
}

func (h *GiveawayHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	var code models.Giveaway

	if err := json.NewDecoder(r.Body).Decode(&code); err != nil {
		log.Printf("Error decoding JSON in CreateCode: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.giveaway.Add(code); err != nil {
		log.Printf("Error adding code: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(code); err != nil {
		log.Printf("Error encoding response in CreateCode: %v", err)
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

func (h *GiveawayHandler) GetRandomCode(w http.ResponseWriter, r *http.Request) {
	code, err := h.giveaway.GetRandomCode()
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoCodes):
			http.Error(w, "We are out of codes.", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Eror getting code: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(code)
	if err != nil {
		http.Error(w, "Something wrong happen", http.StatusBadRequest)
		log.Printf("Encoding json error when getting random code, details: %v", err)
		return
	}
}

func (h *GiveawayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.GetRandomCode(w, r)
	case http.MethodPost:
		h.CreateCode(w, r)
	}
}
