package bestbuy_test

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/couchbase/gocb.v1"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func mustOpenBucket(name string) *gocb.Bucket {
	cluster, err := bestbuy.OpenCluster("couchbase://localhost", "evan", "password")
	if err != nil {
		panic(err)
	}

	bucket, err := bestbuy.OpenBucket(name, cluster)
	if err != nil {
		panic(err)
	}
	return bucket
}

func TestPostSearch(t *testing.T) {
	// start server
	server := bestbuy.Server{}
	go server.Start("8081")
	time.Sleep(1 * time.Second)
	// sleep for one second

	t.Log("Started server on port 8081")
	server.Shutdown()
	t.Log("Closed server on port 8081")

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
	bucket.Close()
}

func TestCounters(t *testing.T) {

	bucket := mustOpenBucket("tests")

	// now scan
	type event struct {
		ProductID string `json:"productId"`
		EpochDay  int    `json:"epochDay"`
		EpochHour int    `json:"epochHour"`
		Count     int    `json:"count"`
	}

	view1 := event{"prod01", 3455, 14, 0}
	// upsert document
	bucket.Upsert("event01", view1, 0)

	// now increment the counter
	bucket.MutateIn("event01", 0, 0).Counter("count", 1, false).Execute()

	var view2 event
	bucket.Get("event01", &view2)

	t.Logf("Value of counter: %d", view2.Count)
	if view2.Count != 1 {
		t.Error("counter field was not incremented!")
	}

	bucket.Close()
}

func TestEpochCalculations(t *testing.T) {
	now := time.Now().Unix()
	day := now / 86400
	hour := now / 3600

	t.Logf("The day is %d and the hour is %d", day, hour)

}

func TestEventID(t *testing.T) {

	event := bestbuy.QueryTrackEvent{}
	event.Count = 0
	event.Hour = 177
	event.Day = 2345
	event.Query = "Red Dead"

	t.Logf("ID: %s", event.ID())

}
