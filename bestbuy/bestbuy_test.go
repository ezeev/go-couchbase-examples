package bestbuy_test

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/couchbase/gocb.v1"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func mustOpenBucket(name string) *gocb.Bucket {
	cluster, err := bestbuy.CbConnect("couchbase://localhost", "evan", "password")
	if err != nil {
		panic(err)
	}

	bucket, err := bestbuy.CbOpenBucket(name, cluster)
	if err != nil {
		panic(err)
	}
	return bucket
}

func TestGetProduct(t *testing.T) {

	/*cluster, err := bestbuy.CbConnect("couchbase://localhost", "evan", "password")
	if err != nil {
		t.Error(err)
	}

	bucket, err := bestbuy.CbOpenBucket("bb-catalog", cluster)
	if err != nil {
		t.Error(err)
	}*/
	bucket := mustOpenBucket("bb-catalog")

	prod1, err := bestbuy.GetProduct("1184298", bucket)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Type attribute of doc returned: %s", prod1.Type)
	if prod1.Type != "Game" {
		t.Error("Type != \"Game\"")
	}

}

func TestBoostSearch(t *testing.T) {
	bucket := mustOpenBucket("bb-catalog")
	res, err := bestbuy.Search("bb-search", "xbox sku:\"3650503\"^1", bucket)
	if err != nil {
		t.Error(err)
	}

	bucket.Close()

	if res.Hits()[0].Id != "3650503" {
		t.Errorf("Expected id 3650503 to be in position. Got %s instead\n", res.Hits()[0])
	}

	fmt.Printf("\tFound %d results in %s\n", res.TotalHits(), res.Took().String())
	for _, hit := range res.Hits() {
		fmt.Printf("\t%s : %f\n", hit.Id, hit.Score)
	}
	// facets
	for _, f := range res.Facets() {
		fmt.Println(f.Field)
		for _, v := range f.Terms {
			fmt.Printf("\t%s : %d\n", v.Term, v.Count)
		}
	}

}

func TestFilteredSearch(t *testing.T) {
	bucket := mustOpenBucket("bb-catalog")
	res, err := bestbuy.Search("bb-search", "red dead +platform:\"PS3 Games\"", bucket)
	if err != nil {
		t.Error(err)
	}
	for i, hit := range res.Hits() {
		fmt.Printf("\tHit %d platform: %s\n", i, hit.Fields)
	}
}

func TestCounters(t *testing.T) {

	bucket := mustOpenBucket("test")

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
