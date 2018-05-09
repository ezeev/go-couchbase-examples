package bestbuy_test

import (
	"testing"

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

}
