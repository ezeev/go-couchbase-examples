package bestbuy

import (
	"github.com/gorilla/mux"
	"gopkg.in/couchbase/gocb.v1"
)

type Server struct {
	Router         *mux.Router
	ProductsBucket *gocb.Bucket
	TrackingBucket *gocb.Bucket
}
