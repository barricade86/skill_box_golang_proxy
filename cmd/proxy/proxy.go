package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var counter = 0
var appHostList string
var proxyHostAndPort string
var appHosts []string

func init() {
	flag.StringVar(&appHostList, "app-hosts-list", "", "Set application hosts, separated by commas. Proxy will redirect requests to them")
	flag.StringVar(&proxyHostAndPort, "proxy-host-and-port", "", "Proxy host and port")
	flag.Parse()
	if len(proxyHostAndPort) == 0 {
		panic("Empty proxy host and port value")
	}

	if len(appHostList) == 0 {
		panic("Empty application hosts list")
	}

	appHosts = strings.Split(appHostList, ",")
	if len(appHosts) < 2 {
		panic("2 application hosts should be provided")
	}
}

func main() {
	http.HandleFunc("/", loadBalancer)
	logrus.Fatal(http.ListenAndServe(proxyHostAndPort, nil))
}

func loadBalancer(response http.ResponseWriter, request *http.Request) {
	proxyUrl := getProxyUrl()
	serveReverseProxy(proxyUrl, response, request)
}

func getProxyUrl() string {
	server := appHosts[counter]
	counter++
	counter = counter % len(appHosts)

	return server
}

func serveReverseProxy(target string, response http.ResponseWriter, request *http.Request) {
	parsedUrl, err := url.Parse(target)
	if err != nil {
		logrus.Errorf("Parsing URL %s error:%s", target, err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedUrl)
	proxy.ServeHTTP(response, request)
}
