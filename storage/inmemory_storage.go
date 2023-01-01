package storage

import (
	"fmt"
	"math/rand"
	"sync"
	"webserver/model"
)

type InMemoryStorage struct {
	mu      sync.Mutex
	storage map[uint64]*model.User
}

func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{storage: make(map[uint64]*model.User)}
}

func (ims *InMemoryStorage) Add(user *model.User) (uint64, error) {
	userID := uint64(rand.Int())
	ims.storage[userID] = user

	return userID, nil
}

func (ims *InMemoryStorage) Delete(userID uint64) error {
	_, ok := ims.storage[userID]
	if !ok {
		return fmt.Errorf("user with id %d not exists", userID)
	}
	ims.mu.Lock()
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
	_, ok := ims.storage[userID]
	if !ok {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	ims.mu.Lock()
	ims.storage[userID] = user
	ims.mu.Lock()

	return nil
}
