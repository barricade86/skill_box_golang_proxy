package storage

import (
	"fmt"
	"math/rand"
	"sync"
	"webserver/internal/model"
)

type InMemoryStorage struct {
	mu      sync.Mutex
	storage map[uint64]*model.User
}

func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{storage: make(map[uint64]*model.User)}
}

func (ims *InMemoryStorage) Add(user *model.User) (uint64, error) {
	ims.mu.Lock()
	userID := uint64(rand.Int())
	ims.storage[userID] = user
	defer ims.mu.Unlock()

	return userID, nil
}

func (ims *InMemoryStorage) Delete(userID uint64) error {
	ims.mu.Lock()
	_, ok := ims.storage[userID]
	if !ok {
		return fmt.Errorf("user with id %d not exists", userID)
	}
	delete(ims.storage, userID)
	defer ims.mu.Unlock()
	return nil
}

func (ims *InMemoryStorage) FindByUserId(userID uint64) (*model.User, error) {
	userData, ok := ims.storage[userID]
	if !ok {
		return nil, fmt.Errorf("user with id %d not exists", userID)
	}

	return userData, nil
}

func (ims *InMemoryStorage) Update(userID uint64, user *model.User) error {
	ims.mu.Lock()
	_, ok := ims.storage[userID]
	if !ok {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	ims.storage[userID] = user
	ims.mu.Unlock()

	return nil
}
