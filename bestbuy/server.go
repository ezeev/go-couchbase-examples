package bestbuy

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/couchbase/gocb.v1"
)

type Server struct {
	Router           *mux.Router
	CouchbaseCluster *gocb.Cluster
	ProductsBucket   *gocb.Bucket
	TrackingBucket   *gocb.Bucket
	EventChan        chan Event
	HttpServer       *http.Server
}

func (s *Server) Start(port string) {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/api/search", s.restEndpoint(s.handleSearch()))

	s.HttpServer = &http.Server{
		Handler:      s.Router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// initialize channels
	s.EventChan = make(chan Event, 0)

	// connect to Couchbase buckets
	cluster := MustOpenCluster("couchbase://localhost", "evan", "password")
	pbucket := MustOpenBucket("bb-catalog", cluster)
	tbucket := MustOpenBucket("bb-tracking", cluster)
	s.CouchbaseCluster = cluster
	s.ProductsBucket = pbucket
	s.TrackingBucket = tbucket

	go s.processEvents()

	log.Fatal(s.HttpServer.ListenAndServe())

	// close connections
	s.Shutdown()

}

func (s *Server) Shutdown() {
	s.TrackingBucket.Close()
	s.ProductsBucket.Close()
	s.CouchbaseCluster.Close()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.HttpServer.Shutdown(ctx)
	log.Println("Shutdown server gracefully")

}

func (s *Server) processEvents() {
	for event := range s.EventChan {
		id := event.ID()
		_, err := s.TrackingBucket.Insert(id, event, 0)
		if err != nil {
			//key already exists, increment counter
			s.TrackingBucket.MutateIn(id, 0, 0).Counter("count", 1, false).Execute()
			log.Printf("Incremented counter for id %s", id)
		} else {
			log.Printf("Added counter for id %s", id)
		}
	}
}
