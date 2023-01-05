package main

import (
	"flag"
	"fmt"
	"github.com/FZambia/tarantool"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
	"webserver/api/router"
	usersStorage "webserver/internal/storage"
)

var storageType string
var tarantoolEndPoint string
var tarantoolConnectionTimeout uint
var httpAppEndPoint string
var dataStorage usersStorage.Storage
var tarantoolClient *tarantool.Connection

func init() {
	flag.StringVar(&httpAppEndPoint, "http-endpoint", "", "Webserver endpoint")
	flag.StringVar(&storageType, "storage-type", "", "Storage type:inmemory or tarantool")
	flag.StringVar(&tarantoolEndPoint, "tarantool-endpoint", "", "tarantool EndPoint")
	flag.UintVar(&tarantoolConnectionTimeout, "connection-timeout", 65, "tarantool connection timeout")
	flag.Parse()

	dataStorage = usersStorage.NewStorage()
	if storageType == "tarantool" {
		opts := tarantool.Opts{
			RequestTimeout: 500 * time.Millisecond,
			User:           os.Getenv("TARANTOOL_LOGIN"),
			Password:       os.Getenv("TARANTOOL_PASSWORD"),
			ConnectTimeout: time.Duration(tarantoolConnectionTimeout) * time.Second,
		}
		var err error
		tarantoolClient, err = tarantool.Connect(tarantoolEndPoint, opts)
		if err != nil {
			panic(fmt.Sprintf("Connection refused: %v", err))
		}

		dataStorage = usersStorage.NewTarantoolStorage(tarantoolClient)
	}
}

func main() {
	defer func() { _ = tarantoolClient.Close() }()
	routerService := router.NewRouterService(dataStorage)
	logrus.Fatal(http.ListenAndServe(httpAppEndPoint, routerService.Init()))
}
