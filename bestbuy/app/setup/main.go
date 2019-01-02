package main

import (
	"log"
	"runtime/debug"
	"strings"

	gocb "gopkg.in/couchbase/gocb.v1"
)

func CreateBuckets(mgr *gocb.ClusterManager) {
	settings := &gocb.BucketSettings{}
	settings.Name = "products"
	settings.Quota = 512
	settings.FlushEnabled = true
	err := mgr.InsertBucket(settings)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Println("WARN: products bucket already exists")
		} else {
			panic(err)
		}
	}
	settings = &gocb.BucketSettings{}
	settings.Name = "clicks"
	settings.Quota = 512
	settings.FlushEnabled = true
	err = mgr.InsertBucket(settings)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Println("WARN: clicks bucket already exists")
		} else {
			panic(err)
		}
	}

}

func main() {

	c, err := gocb.Connect("couchbase://52.24.36.56")
	if err != nil {
		log.Println(err)
	}
	c.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("Error authenticating:  %v", err)
	}
	mgr := c.Manager("Administrator", "password")
	CreateBuckets(mgr)

	productsBucket, err := c.OpenBucket("products", "")
	if err != nil {
		panic(err)
	}
	err = LoadProductsFromFile("/Users/evanpease/Development/datasets/product_data.jsonl", productsBucket)
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
	productsBucket.Close()

}
