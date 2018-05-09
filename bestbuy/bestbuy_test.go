package bestbuy_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func TestGetProduct(t *testing.T) {

	cluster, err := bestbuy.CbConnect("couchbase://localhost", "evan", "password")
	if err != nil {
		t.Error(err)
	}

	bucket, err := bestbuy.CbOpenBucket("bb-catalog", cluster)
	if err != nil {
		t.Error(err)
	}

	prod1, err := bestbuy.GetProduct("1184298", bucket)
	t.Logf("Type attribute of doc returned: %s", prod1.Type)
	if prod1.Type != "Game" {
		t.Error("Type != \"Game\"")
	}

	// FTS
	res, err := bestbuy.Search("bb-search", "gears of war", bucket)
	if err != nil {
		t.Error(err)
	}

	bucket.Close()
	cluster.Close()

	for _, hit := range res.Hits() {
		fmt.Printf("%s\n", hit.Id)
	}
	// facets
	for _, f := range res.Facets() {
		fmt.Println(f.Field)
		for k, v := range f.Terms {
			fmt.Printf("\t%s : %s\n", k, v)
		}
	}

}

func TestCounters(t *testing.T) {

	cluster, err := bestbuy.CbConnect("couchbase://localhost", "evan", "password")
	if err != nil {
		t.Error(err)
	}

	bucket, err := bestbuy.CbOpenBucket("tests", cluster)
	if err != nil {
		t.Error(err)
	}

	// use KV
	rows := 10
	i := 0
	for i < rows {
		key := fmt.Sprintf("PROD%d-%d", i, time.Now().Unix())
		bucket.Insert(key, 0, 0)
		i++
	}

	// now increment
	i = 0
	for i < rows {
		key := fmt.Sprintf("PROD%d-%d", i, time.Now().Unix())
		bucket.Counter(key, 1, 0, 0)
	}

	// now scan

}
