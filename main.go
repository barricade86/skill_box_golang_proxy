package main

import (
	//"context"
	"flag"
	"fmt"
	"github.com/FZambia/tarantool"
	"net/http"
	"os"
	"time"
	"webserver/controller"
	"webserver/storage"
)

var storageType string
var tarantoolEndPoint string
var tarantoolConnectionTimeout uint
var httpEndPoint string
var dataStorage storage.Storage

func init() {
	flag.StringVar(&httpEndPoint, "http-endpoint", "", "Webserver endpoint")
	flag.StringVar(&storageType, "storage-type", "", "Storage type:inmemory or tarantool")
	flag.StringVar(&tarantoolEndPoint, "tarantool-endpoint", "", "tarantool EndPoint")
	flag.UintVar(&tarantoolConnectionTimeout, "connection-timeout", 5, "tarantool connection timeout")
	flag.Parse()
}

func main() {
	if storageType == "tarantool" {
		opts := tarantool.Opts{
			RequestTimeout: 500 * time.Millisecond,
			User:           os.Getenv("TARANTOOL_LOGIN"),
			Password:       os.Getenv("TARANTOOL_PASSWORD"),
			ConnectTimeout: time.Duration(tarantoolConnectionTimeout) * time.Second,
		}
		client, err := tarantool.Connect(tarantoolEndPoint, opts)
		if err != nil {
			panic(fmt.Sprintf("Connection refused: %v", err))
		}

		defer func() { _ = client.Close() }()

		dataStorage = storage.NewTarantoolStorage(client)
	} else {
		dataStorage = storage.NewStorage()
	}

	ctrl := controller.NewController(dataStorage)
	mux := http.NewServeMux()
	mux.HandleFunc("/create", ctrl.Create)
	mux.HandleFunc("/make_friends", ctrl.MakeFriends)
	mux.HandleFunc("/delete", ctrl.Delete)
	mux.HandleFunc("/friends/", ctrl.GetFriendsByUserId)
	mux.HandleFunc("/", ctrl.UpdateAgeByUserId)

	http.ListenAndServe(httpEndPoint, mux)
}
