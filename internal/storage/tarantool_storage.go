package storage

import (
	"fmt"
	"github.com/FZambia/tarantool"
	"reflect"
	"webserver/internal/model"
)

type userResponse struct {
	Id      interface{}
	Name    string
	Age     int
	Friends []uint64
}

// TarantoolStorage Struct for tarantool storage
type TarantoolStorage struct {
	client *tarantool.Connection
}

// NewTarantoolStorage Creates the instance for MongoDB
func NewTarantoolStorage(tarantoolClient *tarantool.Connection) *TarantoolStorage {
	return &TarantoolStorage{
		client: tarantoolClient,
	}
}

// Add insert new record into collection
func (t *TarantoolStorage) Add(user *model.User) (uint64, error) {
	result, err := t.client.Exec(tarantool.Insert("users", &userResponse{nil, user.Name, user.Age, user.Friends}))
	if err != nil {
		return 0, fmt.Errorf("insert data record error:", err)
	}

	userID := reflect.ValueOf(result[0]).Index(0).Interface().(int64)

	return uint64(userID), nil
}

// Delete Deletes record by id
func (t *TarantoolStorage) Delete(userID uint64) error {
	_, err := t.client.Exec(tarantool.Delete("users", 0, []interface{}{userID}))
	if err != nil {
		return fmt.Errorf("delete user error:%s", err)
	}

	return nil
}

// FindByUserId find information by userID
func (t *TarantoolStorage) FindByUserId(userID uint64) (*model.User, error) {
	var users []model.User
	err := t.client.ExecTyped(tarantool.Select("users", "primary", 0, 100, tarantool.IterEq, tarantool.UintKey{userID}), &users)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("The result is empty")
	}

	return &users[0], nil
}

// Update updates info about user
func (t *TarantoolStorage) Update(userID uint64, user *model.User) error {
	_, err := t.client.Exec(
		tarantool.Update(
			"users",
			"primary",
			[]interface{}{userID},
			[]tarantool.Op{
				tarantool.OpAssign(1, user.Id),
				tarantool.OpAssign(2, user.Name),
				tarantool.OpAssign(3, user.Age),
				tarantool.OpAssign(4, user.Friends),
			},
		))
	if err != nil {
		return err
	}

	return nil
}
