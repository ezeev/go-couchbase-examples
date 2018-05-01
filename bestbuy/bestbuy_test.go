package bestbuy_test

import (
	"testing"

	"github.com/ezeev/go-couchbase-examples/bestbuy"
)

func TestCatalogSave(t *testing.T) {
	err := bestbuy.LoadProducts("data/products.jsonl")
	if err != nil {
		t.Error(err)
	}
}
