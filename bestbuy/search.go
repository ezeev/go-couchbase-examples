package bestbuy

import (
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/couchbase/gocb.v1/cbft"
)

func Search(index string, q string, bucket *gocb.Bucket) (gocb.SearchResults, error) {
	//query := gocb.NewSearchQuery(index, cbft.NewMatchQuery(q))
	query := gocb.NewSearchQuery(index, cbft.NewQueryStringQuery(q))

	query.AddFacet("manufacturer", cbft.NewTermFacet("manufacturer", 10))
	query.AddFacet("platform", cbft.NewTermFacet("platform", 10))
	query.Fields("platform")
	res, err := bucket.ExecuteSearchQuery(query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
