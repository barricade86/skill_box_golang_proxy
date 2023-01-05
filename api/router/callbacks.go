package router

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"webserver/internal/model"
	requestpkg "webserver/internal/request"
)

// create Creates new user
func (rs *RouterService) create(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Status Method not allowed`))
		return
	}

	jsonBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logrus.Errorf("Read request body error:%s", err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	userData := &requestpkg.UserData{}
	err = json.Unmarshal(jsonBody, &userData)
	if err != nil {
		logrus.Errorf("Status internal server error due to unmarshal body:%s", err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("Status internal server error due to unmarshal body:%s", err)))
		return
	}

	user := &model.User{Name: userData.Name, Age: userData.Age, Friends: userData.Friends}
	userID, err := rs.dataStorage.Add(user)
	if err != nil {
		logrus.Errorf("Storage addition error:%s", err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	logrus.Infof("User with id %d was created", userID)
	response.Write([]byte(fmt.Sprintf(`{"userID":%d}`, userID)))
	response.WriteHeader(http.StatusCreated)
}

// addFriendsForUser Add Friends For User
func (rs *RouterService) addFriendsForUser(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Status Bad Request`))
		return
	}

	jsonBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	friendshipRequest := &requestpkg.FriendshipRequest{}
	err = json.Unmarshal(jsonBody, &friendshipRequest)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user, err := rs.dataStorage.FindByUserId(friendshipRequest.SourceId)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	if user == nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`user with requested ID not found`))
		return
	}

	user.Friends = append(user.Friends, friendshipRequest.TargetId)
	secondUser, err := rs.dataStorage.FindByUserId(friendshipRequest.TargetId)
	err = rs.dataStorage.Update(friendshipRequest.SourceId, user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	err = rs.dataStorage.Update(friendshipRequest.TargetId, secondUser)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	response.Write([]byte(user.Name + " and " + secondUser.Name + " are friends now"))
	response.WriteHeader(http.StatusOK)
}

// delete Deletes user
func (rs *RouterService) delete(response http.ResponseWriter, request *http.Request) {
	if request.Method != "DELETE" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Status Bad Request`))
		return
	}

	jsonBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	userDelete := &requestpkg.UserDelete{}
	err = json.Unmarshal(jsonBody, &userDelete)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user, err := rs.dataStorage.FindByUserId(userDelete.TargetId)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	if user == nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`user with requested ID not found`))
		return
	}

	response.Write([]byte(user.Name + ` was deleted`))
	response.WriteHeader(http.StatusOK)
}

// getFriendsForUser shows friends of user
func (rs *RouterService) getFriendsForUser(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Status Bad Request`))
		return
	}

	userID, err := strconv.ParseUint(chi.URLParam(request, "userID"), 10, 64)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("Error parsing user data. userID PathParameter is not valid:%s", err)))
		return
	}

	user, err := rs.dataStorage.FindByUserId(userID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	if user == nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(fmt.Sprintf("user with requested id %d not found", userID)))
		return
	}

	for _, val := range user.Friends {
		friend, fbUidErr := rs.dataStorage.FindByUserId(val)
		if fbUidErr != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", fbUidErr)))
			return
		}

		response.Write([]byte(friend.Name))
	}

	response.WriteHeader(http.StatusOK)
}

// updateAgeByUserId Updates age of user by userID
func (rs *RouterService) updateAgeByUserId(response http.ResponseWriter, request *http.Request) {
	if request.Method != "PUT" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Status Bad Request`))
		return
	}

	userID, err := strconv.ParseUint(chi.URLParam(request, "userID"), 10, 64)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("Error parsing user data. userID PathParameter is not valid:%s", err)))
		return
	}

	jsonBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	ageUpdate := &requestpkg.AgeUpdate{}
	err = json.Unmarshal(jsonBody, &ageUpdate)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user, err := rs.dataStorage.FindByUserId(userID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	if user == nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(fmt.Sprintf("user with requested id %d not found", userID)))
		return
	}

	user.Age = ageUpdate.Age
	_, err = rs.dataStorage.Add(user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Status Internal Error` + fmt.Sprintf("%s", err)))
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write([]byte(fmt.Sprintf("age of %s successfully updated. Now his age is %d", user.Name, user.Age)))
}
