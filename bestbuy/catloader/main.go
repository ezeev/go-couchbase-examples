package main

import (
	"log"
	"os"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func main() {

	args := os.Args
	path := args[1]

	//err := bestbuy.LoadProductsFromFile("/Users/evanpease/Development/datasets/product_data.jsonl")
	err := bestbuy.LoadProductsFromFile(path)
	if err != nil {
		log.Println(err)
	}

}
