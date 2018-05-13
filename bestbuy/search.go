package bestbuy

import (
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/couchbase/gocb.v1/cbft"
)

func Search(index string, q string, bucket *gocb.Bucket) (SearchResponse, error) {
	//query := gocb.NewSearchQuery(index, cbft.NewMatchQuery(q))
	query := gocb.NewSearchQuery(index, cbft.NewQueryStringQuery(q))

	query.AddFacet("manufacturer", cbft.NewTermFacet("manufacturer", 10))
	query.AddFacet("platform", cbft.NewTermFacet("platform", 10))
	query.Fields("platform")
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
	return searchResponse, nil
}

type SearchResponse struct {
	Query     string                            `json:"query"`
	Hits      []gocb.SearchResultHit            `json:"hits"`
	Facets    map[string]gocb.SearchResultFacet `json:"facets"`
	Errors    []string                          `json:"errors"`
	Status    gocb.SearchResultStatus           `json:"status"`
	TotalHits int                               `json:"totalHits"`
	MaxScore  float64                           `json:"maxScore"`
	TookMs    float64                           `json:"tookMs"`
}
