package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/ezeev/go-couchbase-examples/bestbuy/app/model"
	gocb "gopkg.in/couchbase/gocb.v1"
)

func LoadProductsFromFile(pathToJson string, bucket *gocb.Bucket) error {

	file, err := ioutil.ReadFile(pathToJson)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(file)
	for {
		s, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		product := model.Product{}
		err = json.Unmarshal([]byte(s), &product)
		if err != nil {
			return err
		}
		err = SaveProduct(product, bucket)
		if err != nil {
			return err
		}
	}
	err = bucket.Close()
	return err
}

func SaveProduct(product model.Product, bucket *gocb.Bucket) error {
	_, err := bucket.Upsert(product.ID, product, 0)
	return err
}
