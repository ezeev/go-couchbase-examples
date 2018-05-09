package main

import (
	"log"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func main() {

	err := bestbuy.LoadProductsFromFile("../data/products.jsonl")
	if err != nil {
		log.Println(err)
	}

}
