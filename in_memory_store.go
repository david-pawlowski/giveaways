package main

import "errors"

var ErrNotFound = errors.New("Could not find code with given id")

type InMemoryStore struct {
	store map[int]Code
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{map[int]Code{1: {"abC", false}}}
}

func (ms *InMemoryStore) Add(id int, code Code) {
	ms.store[id] = code
}

func (ms *InMemoryStore) Get(id int) (Code, error) {
	code, ok := ms.store[id]
	if !ok {
		return code, ErrNotFound
	}
	return code, nil
}

func (ms *InMemoryStore) MarkClaimed(id int) bool {
	code := ms.store[id]
	code.Claimed = true
	ms.store[id] = code
	return true
}
