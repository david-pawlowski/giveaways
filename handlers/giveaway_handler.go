package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/david-pawlowski/giveaway/models"
	"github.com/david-pawlowski/giveaway/repository"
	"github.com/david-pawlowski/giveaway/service"
	"net/http"
	"strconv"
	"strings"
)

type GiveawayHandler struct {
	giveaway service.GiveawayService
}

func NewGiveawayHandler(g service.GiveawayService) *GiveawayHandler {
	return &GiveawayHandler{
		giveaway: g,
	}
}

func (gh *GiveawayHandler) GetCode(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.EscapedPath(), "/")
	id, err := strconv.Atoi(path[len(path)-1])
	data, err := gh.giveaway.Get(id)
	if err == repository.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found item for given id"))
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(data)
	if err != nil {
		fmt.Print(err)
	}
	gh.giveaway.MarkClaimed(id)
	w.Write(jsonResp)
	return
}

func (gh *GiveawayHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	var code models.Code

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	code, err = gh.giveaway.Add(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Can't create code with this id"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(code)
}

func (gh *GiveawayHandler) GetRandomCode(w http.ResponseWriter, r *http.Request) {
	code, err := gh.giveaway.GetRandomCode()
	if err == repository.ErrNoCodes {
		w.Write([]byte("No more codes available"))
		return
	}
	if err != nil {
		http.Error(w, "Something wrong happen", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)

	jsonResp, err := json.Marshal(code)
	if err != nil {
		fmt.Print(err)
	}
	w.Write(jsonResp)
}

func (gh *GiveawayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	switch r.Method {
	case http.MethodGet:
		gh.GetRandomCode(w, r)
	case http.MethodPost:
		gh.CreateCode(w, r)
	}
}
