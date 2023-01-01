package controller

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"webserver/model"
	"webserver/request"
	"webserver/storage"
)

type Controller struct {
	storage storage.Storage
}

func NewController(storage storage.Storage) *Controller {
	return &Controller{storage: storage}
}

func (c *Controller) Create(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Read request body error:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	userData := &request.UserData{}
	err = json.Unmarshal(jsonBody, &userData)
	if err != nil {
		logrus.Error("Status internal server error due to unmarshal body:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf("Status internal server error due to unmarshal body:%s", err)))
		return
	}

	user := &model.User{Name: userData.Name, Age: userData.Age, Friends: userData.Friends}
	userID, err := c.storage.Add(user)
	if err != nil {
		logrus.Error("Storage addition error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	logrus.Infof("User with id %d was created", userID)
	rw.Write([]byte(fmt.Sprintf(`{"userID":%d}`, userID)))
	rw.WriteHeader(http.StatusCreated)
}

func (c *Controller) MakeFriends(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	friendshipRequest := &request.FriendshipRequest{}
	err = json.Unmarshal(jsonBody, &friendshipRequest)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user, err := c.storage.FindByUserId(friendshipRequest.SourceId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	if user == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`user with requested ID not found`))
		return
	}

	user.Friends = append(user.Friends, friendshipRequest.TargetId)
	secondUser, err := c.storage.FindByUserId(friendshipRequest.TargetId)
	err = c.storage.Update(friendshipRequest.SourceId, user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	err = c.storage.Update(friendshipRequest.TargetId, secondUser)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	rw.Write([]byte(user.Name + " and " + secondUser.Name + " are friends now"))
	rw.WriteHeader(http.StatusOK)
}

func (c *Controller) Delete(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	userDelete := &request.UserDelete{}
	err = json.Unmarshal(jsonBody, &userDelete)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user, err := c.storage.FindByUserId(userDelete.TargetId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	if user == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`user with requested ID not found`))
		return
	}

	rw.Write([]byte(user.Name + ` was deleted`))
	rw.WriteHeader(http.StatusOK)
}

func (c *Controller) GetFriendsByUserId(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	userID, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/friends/"), 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Error parsing user data`))
		return
	}

	user, err := c.storage.FindByUserId(userID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	if user == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("user with requested id %d not found", userID)))
		return
	}

	for _, val := range user.Friends {
		friend, fbUidErr := c.storage.FindByUserId(val)
		if fbUidErr != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", fbUidErr)))
			return
		}

		rw.Write([]byte(friend.Name))
	}

	rw.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateAgeByUserId(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	userID, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/"), 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Error parsing user data`))
		return
	}

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	ageUpdate := &request.AgeUpdate{}
	err = json.Unmarshal(jsonBody, &ageUpdate)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user, err := c.storage.FindByUserId(userID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	if user == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("user with requested id %d not found", userID)))
		return
	}

	user.Age = ageUpdate.Age
	_, err = c.storage.Add(user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("возраст пользователя %s успешно обновлён. Теперь ему %d", user.Name, user.Age)))
}
