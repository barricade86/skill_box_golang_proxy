package main

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var counter = 0
var firstInstanceHost = "http://localhost:8080"
var secondInstanceHost = "http://localhost:8081"

func main() {
	http.HandleFunc("/create", handle)
	http.HandleFunc("/make_friends", handle)
	http.HandleFunc("/delete", handle)
	http.HandleFunc("/friends", handle)
	http.HandleFunc("/", handle)
	http.ListenAndServe("localhost:9000", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Error reading raw body data")
		return
	}

	defer r.Body.Close()
	requestURI := ""
	if counter == 0 {
		requestURI = firstInstanceHost + r.RequestURI
		counter++
	} else {
		requestURI = secondInstanceHost + r.RequestURI
		counter--
	}

	var response *http.Response
	switch r.Method {
	case "POST":
		response, err = http.Post(requestURI, "application/json", r.Body)
	case "GET":
		response, err = http.Get(requestURI)
	case "PUT":
		putRequest, err := http.NewRequest("PUT", requestURI, r.Body)
		if err != nil {
			logrus.Error("put request creation error:", err)
			return
		}

		response, err = http.DefaultClient.Do(putRequest)
	case "DELETE":
		deleteRequest, err := http.NewRequest("PUT", requestURI, r.Body)
		if err != nil {
			logrus.Error("request creation error:", err)
			return
		}

		response, err = http.DefaultClient.Do(deleteRequest)
	}

	if err != nil {
		logrus.Error("send delete request error:", err)
		return
	}

	logrus.Infof("response data:%v", response.Body)
}
