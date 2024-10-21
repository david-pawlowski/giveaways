package handlers

import (
	"encoding/json"
	"github.com/david-pawlowski/giveaway/models"
	"github.com/david-pawlowski/giveaway/repository"
	"github.com/david-pawlowski/giveaway/service"
	"net/http"
)

type GiveawayHandler struct {
	giveaway service.GiveawayService
}

func NewGiveawayHandler(g service.GiveawayService) *GiveawayHandler {
	return &GiveawayHandler{
		giveaway: g,
	}
}

func (gh *GiveawayHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	var code models.Giveaway

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	gh.giveaway.Add(code)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(code)
}

func (gh *GiveawayHandler) GetRandomCode(w http.ResponseWriter, r *http.Request) {
	code, err := gh.giveaway.GetRandomCode()
	if err == repository.ErrNoCodes {
		http.Error(w, "No more codes left", http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(code)
	if err != nil {
		http.Error(w, "Something wrong happen", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func (gh *GiveawayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		gh.GetRandomCode(w, r)
	case http.MethodPost:
		gh.CreateCode(w, r)
	}
}
