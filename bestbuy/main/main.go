package main

import (
	"os"
	"os/signal"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	bbserver := bestbuy.Server{}
	go bbserver.Start("8081")

	<-stop
	bbserver.Shutdown()

}
