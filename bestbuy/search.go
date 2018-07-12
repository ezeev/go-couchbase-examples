package bestbuy

import (
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/couchbase/gocb.v1/cbft"
)

func CountField(key string, bucket *gocb.Bucket) int64 {
	var event ClickTrackEvent
	bucket.Get(key, &event)
	return event.Count
}

func SignalSearch(index string, q string, numSignals int, bucket *gocb.Bucket) (map[string]int64, error) {

	matchQuery := cbft.NewMatchQuery(q)
	matchQuery.Analyzer("en")
	query := gocb.NewSearchQuery(index, cbft.NewDisjunctionQuery(matchQuery))
	query.Limit(numSignals)
	query.Fields("*")
	res, err := bucket.ExecuteSearchQuery(query)
	if err != nil {
		return nil, err
	}

	aggSigs := make(map[string]int64)
	for _, v := range res.Hits() {
		//fmt.Printf("%s : %d\n", v.Fields["sku"], CountField(v.Id, bucket))
		aggSigs[v.Fields["sku"]] += CountField(v.Id, bucket)
	}
	return aggSigs, nil

}

func Search(index string, q interface{}, bucket *gocb.Bucket) (SearchResponse, error) {

	query := gocb.NewSearchQuery(index, q)

	query.AddFacet("manufacturer", cbft.NewTermFacet("manufacturer", 10))
	query.AddFacet("platform", cbft.NewTermFacet("platform", 10))
	query.Fields("platform")
	query.Limit(12)

	res, err := bucket.ExecuteSearchQuery(query)
	if err != nil {
		return SearchResponse{}, err
	}
	searchResponse := SearchResponse{
		Query:     q,
		Hits:      res.Hits(),
		Facets:    res.Facets(),
		Errors:    res.Errors(),
		Status:    res.Status(),
		TotalHits: res.TotalHits(),
		MaxScore:  res.MaxScore(),
		TookMs:    res.Took().Seconds() * 1000,
	}

	//add product docs
	searchResponse.Products = GetProductsFromHits(res.Hits(), bucket)
	return searchResponse, nil
}

type SearchResponse struct {
	Q         string                            `json:"q"`
	Session   string                            `json:"session"`
	Query     interface{}                       `json:"query"`
	Hits      []gocb.SearchResultHit            `json:"hits"`
	Products  []*Product                        `json:"products"`
	Facets    map[string]gocb.SearchResultFacet `json:"facets"`
	Errors    []string                          `json:"errors"`
	Status    gocb.SearchResultStatus           `json:"status"`
	TotalHits int                               `json:"totalHits"`
	MaxScore  float64                           `json:"maxScore"`
	TookMs    float64                           `json:"tookMs"`
}
