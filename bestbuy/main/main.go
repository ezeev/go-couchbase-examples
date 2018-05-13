package main

import (
	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func main() {
	bbserver := bestbuy.Server{}
	bbserver.Start("8081")
}
