package storage

import (
	"webserver/internal/model"
)

type Storage interface {
	Add(user *model.User) (uint64, error)
	Delete(userID uint64) error
	FindByUserId(userID uint64) (*model.User, error)
	Update(userID uint64, user *model.User) error
}
