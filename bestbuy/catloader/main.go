package main

import (
	"log"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func main() {

	err := bestbuy.LoadProductsFromFile("/Users/evanpease/Development/datasets/product_data.jsonl")
	if err != nil {
		log.Println(err)
	}

}
