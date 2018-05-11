package main

import "github.com/ezeev/go-couchbase-examples/bestbuy"

func main() {

	cluster, err := bestbuy.CbConnect("couchbase://localhost", "evan", "password")
	if err != nil {
		panic(err)
	}

	pbucket, err := bestbuy.CbOpenBucket("bb-catalog", cluster)
	if err != nil {
		panic(err)
	}

	tbucket, err := bestbuy.CbOpenBucket("bb-tracking", cluster)
	if err != nil {
		panic(err)
	}
	bbserver := bestbuy.Server{
		ProductsBucket: pbucket,
		TrackingBucket: tbucket,
	}
	bbserver.Start("8081")

}
