package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	conf := bestbuy.Config{
		APIPort:                 8081,
		CouchbaseUserName:       "admin",
		CouchbasePassword:       "password123",
		CouchbaseClusterAddress: "couchbase://localhost",
		AssetsBaseUrl:           "http://localhost:8082",
		ProductsFTSIndexName:    "productsfts",
		TrackingFTSIndexName:    "clicksfts",
	}

	bbserver := bestbuy.NewServer(&conf)
	go bbserver.Start()

	// static file server
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	go http.ListenAndServe(":8082", nil)

	<-stop
	bbserver.Shutdown()

}
