package bestbuy

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	gocb "gopkg.in/couchbase/gocb.v1"
	"gopkg.in/couchbase/gocb.v1/cbft"
)

func httpClient() *http.Client {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	return netClient
}

func configForTest() *Config {
	conf := Config{
		APIPort:                 8081,
		CouchbaseUserName:       "admin",
		CouchbasePassword:       "password123",
		CouchbaseClusterAddress: "couchbase://localhost",
		AssetsBaseUrl:           "http://localhost:8082",
		ProductsFTSIndexName:    "productsfts",
		TrackingFTSIndexName:    "clicksfts",
	}
	return &conf
}

func mustOpenBucket(name string) *gocb.Bucket {
	cluster, err := OpenCluster("couchbase://localhost", "admin", "password123")
	if err != nil {
		panic(err)
	}

	bucket, err := OpenBucket(name, cluster)
	if err != nil {
		panic(err)
	}
	return bucket
}

func TestProductAPI(t *testing.T) {
	server := NewServer(configForTest())
	go server.Start()
	server.ServerReady()

	cli := httpClient()

	url := "http://127.0.0.1:8081/api/product/1000026"

	resp, err := cli.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		t.Error(err)
	}
	if product.Platform == "Comedy/Spoken" {
		t.Log("Sku ID has platform: Comedy/Spoken")
	} else {
		t.Errorf("Expected platform = Comedy/Spoken")
	}

	server.Shutdown()

}

func TestGetProducts(t *testing.T) {
	prodIds := []string{"3725313", "9891968", "4979434", "4979407", "4979073", "5470038", "5470056", "2319054", "3504297", "3086761"}
	bucket := mustOpenBucket("products")
	res := GetProducts(prodIds, bucket)
	for _, v := range res {
		if v.Name == "" {
			t.Error("Name field is empty?")
		}
	}
}

func TestSearchTracking(t *testing.T) {
	// start server
	server := NewServer(configForTest())
	go server.Start()
	server.ServerReady()

	// run tests here
	cli := httpClient()

	// 1. Execute a search
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8081/api/search", nil)
	q := req.URL.Query()
	q.Add("track", "true")
	q.Add("q", "TESTING")
	q.Add("session", "TEST")
	req.URL.RawQuery = q.Encode()

	resp, err := cli.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	searchResp := SearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&searchResp)
	if err != nil {
		t.Error(err)
	}
	for _, v := range searchResp.Hits {
		t.Logf("Hit id: %s\n", v.Id)
	}

	server.Shutdown()

	// now run a n1ql query to see if search track even is present.
	bucket := mustOpenBucket("tracking")
	query := gocb.NewN1qlQuery("select * from tracking where query = 'TESTING'")
	rows, err := bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		t.Error(err)
	}
	var row interface{}
	if !rows.Next(&row) {
		t.Error("Search track event was not added!")
	}

	// was the session added?
	sessBucket := mustOpenBucket("sessions")
	var sess Session
	_, err = sessBucket.Get("TEST", sess)
	if err != nil {
		t.Error(err)
	}

	// now clean up
	_, err = sessBucket.Remove("TEST", 0)
	if err != nil {
		t.Error(err)
	}
	// clean up
	query = gocb.NewN1qlQuery("delete from tracking where query = 'TESTING'")
	_, err = bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		t.Error(err)
	}

}

func TestClickTracking(t *testing.T) {
	server := NewServer(configForTest())
	go server.Start()
	server.ServerReady()

	//click a fakeproduct
	client := httpClient()

	// /api/product/<prod ID>?track=true&q=<referring search>
	url := "http://127.0.0.1:8081/api/product/4905704"
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("track", "true")
	q.Add("q", "test click through")
	q.Add("session", "TEST")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.Status)
	server.Shutdown()

	// now make sure the event is there
	bucket := mustOpenBucket("tracking")
	query := gocb.NewN1qlQuery("select * from tracking where query = 'test click through'")
	rows, err := bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		t.Error(err)
	}
	var row interface{}
	if !rows.Next(&row) {
		t.Error("Search track event was not added!")
	}

	// did the session get created?
	sessBucket := mustOpenBucket("sessions")
	var sess Session
	_, err = sessBucket.Get("TEST", &sess)
	if err != nil {
		t.Error("TEST session does not exist!")
	}
	if sess.ProductViews["4905704"] != 1 {
		t.Error("Expected 1 view for product ID 4905704 in session TEST")
	}

	// clean up
	query = gocb.NewN1qlQuery("delete from tracking where query = 'test click through'")
	_, err = bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		t.Error(err)
	}

	_, err = sessBucket.Remove("TEST", 0)
	if err != nil {
		t.Error(err)
	}

}

func TestGetProduct(t *testing.T) {

	bucket := mustOpenBucket("products")

	prod1, err := GetProduct("1184298", bucket)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Type attribute of doc returned: %s", prod1.Type)
	if prod1.Type != "product" {
		t.Error("Type != \"product\"")
	}

}

func TestSignalSearch(t *testing.T) {
	bucket := mustOpenBucket("tracking")
	defer bucket.Close()
	q := "ipad"

	sigs, err := SignalSearch("clicksfts", q, 100, bucket)
	if err != nil {
		t.Error(err)
	}
	for k, v := range sigs {
		v := float64(v)
		v = (math.Log(v) / 10) + 0.01
		t.Logf("%s : %f", k, v)
	}

}

func TestBoostSearch(t *testing.T) {
	bucket := mustOpenBucket("products")

	q := cbft.NewMatchQuery("xbox")
	bq := cbft.NewMatchQuery("3650503")
	bq.Field("sku")
	bq.Boost(1.0)

	res, err := Search("productsfts", cbft.NewDisjunctionQuery(q, bq), bucket)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	bucket.Close()

	if len(res.Hits) == 0 {
		t.Error("No hits!")
		t.Fail()
	}

	//if res.Hits()[0].Id != "3650503" {
	if res.Hits[0].Id != "3650503" {
		t.Errorf("Expected id 3650503 to be in position. Got %s instead\n", res.Hits[0].Id)
	}

	t.Logf("\tFound %d results in %f\n", res.TotalHits, res.TookMs)
	for _, hit := range res.Hits {
		fmt.Printf("\t%s : %f\n", hit.Id, hit.Score)
	}
	// facets
	for _, f := range res.Facets {
		fmt.Println(f.Field)
		for _, v := range f.Terms {
			fmt.Printf("\t%s : %d\n", v.Term, v.Count)
		}
	}

}

func TestFilteredSearch(t *testing.T) {
	bucket := mustOpenBucket("products")
	res, err := Search("productsfts", cbft.NewQueryStringQuery("red dead +platform:\"PS3 Games\""), bucket)
	if err != nil {
		t.Error(err)
	}
	for i, hit := range res.Hits {
		t.Logf("\nHit %d ID: %s\nName: %s\n", i, hit.Id, res.Products[i].Name)
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

	// increment again (why not?)
	bucket.MutateIn("event01", 0, 0).Counter("count", 1, false).Execute()

	bucket.Close()
}

func TestEpochCalculations(t *testing.T) {
	now := time.Now().Unix()
	day := now / 86400
	hour := now / 3600

	t.Logf("The day is %d and the hour is %d", day, hour)

}

func TestEventID(t *testing.T) {

	event := QueryTrackEvent{}
	event.Count = 0
	event.Hour = 177
	event.Day = 2345
	event.Query = "Red Dead"
	t.Logf("ID: %s", event.ID())
}

// BENCHMARKS

func Benchmark10Products(b *testing.B) {
	prodIds := []string{"3725313", "9891968", "4979434", "4979407", "4979073", "5470038", "5470056", "2319054", "3504297", "3086761"}
	products := make([]Product, len(prodIds))

	bucket := mustOpenBucket("products")
	for i, key := range prodIds {
		var prod Product
		bucket.Get(key, &prod)
		products[i] = prod
	}
	bucket.Close()
}

func BenchmarkDocCounter(b *testing.B) {
	bucket := mustOpenBucket("tests")

	type event struct {
		ProductID string `json:"productId"`
		EpochDay  int    `json:"epochDay"`
		EpochHour int    `json:"epochHour"`
		Count     int    `json:"count"`
	}

	view1 := event{"prod01", 3455, 14, 0}
	bucket.Upsert("event01", view1, 0)
	// now increment the counter
	limit := 1000
	b.StartTimer()
	for i := 0; i < limit; i++ {
		bucket.MutateIn("event01", 0, 0).Counter("count", 1, false).Execute()
	}
	b.StopTimer()
	bucket.Close()
}

func BenchmarkInserts(b *testing.B) {
	bucket := mustOpenBucket("tests")

	type event struct {
		ProductID string `json:"productId"`
		EpochDay  int    `json:"epochDay"`
		EpochHour int    `json:"epochHour"`
	}

	// do 100 inserts
	limit := 1000
	b.StartTimer()
	for i := 0; i < limit; i++ {
		bucket.Insert(fmt.Sprintf("test-id-%d", i), event{ProductID: "test", EpochDay: 3444, EpochHour: 244444}, 0)
	}
	b.StopTimer()
	bucket.Close()

}

func BenchmarkKVCounter(b *testing.B) {
	bucket := mustOpenBucket("tests")
	bucket.Upsert("test-counter", 0, 0)

	limit := 1000
	b.StartTimer()
	for i := 0; i < limit; i++ {
		bucket.Counter("test-counter", 1, 0, 0)
	}
	b.StopTimer()
	bucket.Close()

}
