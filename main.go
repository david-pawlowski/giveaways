package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	store := NewInMemoryStore()
	giveawayHandler := NewGiveawayHandler(store)
	mux := http.NewServeMux()
	mux.Handle("/", &homeHandler{})
	mux.Handle("/codes/", giveawayHandler)
	mux.Handle("/codes", giveawayHandler)
	http.ListenAndServe(":8080", mux)
}

type Code struct {
	Code    string `json:"code"`
	Claimed bool   `json:"claimed"`
}

type Giveaway interface {
	Add(code Code) (Code, error)
	Get(id int) (Code, error)
	MarkClaimed(id int) bool
}

type GiveawayHandler struct {
	giveaway Giveaway
}

func NewGiveawayHandler(g Giveaway) *GiveawayHandler {
	return &GiveawayHandler{
		giveaway: g,
	}
}

func (gh *GiveawayHandler) GetCode(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.EscapedPath(), "/")
	id, err := strconv.Atoi(path[len(path)-1])
	data, err := gh.giveaway.Get(id)
	if err == ErrNotFound {
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
	var code Code

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	code, err = gh.giveaway.Add(code)
	fmt.Print(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Can't create code with this id"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(code)
}

func (gh *GiveawayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		gh.GetCode(w, r)
	case http.MethodPost:
		gh.CreateCode(w, r)
	}
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
