package main

import "errors"

var (
	ErrNotFound      = errors.New("Could not find code with given id")
	ErrAlreadyExists = errors.New("Code with this id already exists")
	ErrNoCodes       = errors.New("There are no more codes left")
)

type InMemoryStore struct {
	store map[int]Code
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{map[int]Code{0: {"abC", false}}}
}

func (ms *InMemoryStore) Add(code Code) (Code, error) {
	id := len(ms.store) + 1
	_, err := ms.Get(id)
	if err != nil {
		ms.store[id] = code
		return code, nil
	} else {
		return code, ErrAlreadyExists
	}
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

func (ms *InMemoryStore) GetRandomCode() (Code, error) {
	for i := 0; i < len(ms.store); i++ {
		if !ms.store[i].Claimed {
			ms.MarkClaimed(i)
			return ms.store[i], nil
		}
	}
	return Code{}, ErrNoCodes
}
