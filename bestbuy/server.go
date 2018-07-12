package bestbuy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/couchbase/gocb.v1"
)

type Server struct {
	Router           *mux.Router
	Port             int
	CouchbaseCluster *gocb.Cluster
	ProductsBucket   *gocb.Bucket
	TrackingBucket   *gocb.Bucket
	SessionsBucket   *gocb.Bucket
	EventChan        chan Event
	HttpServer       *http.Server
	Config           *Config
}

func NewServer(config *Config) *Server {

	s := Server{}
	s.Config = config
	s.Port = config.APIPort
	s.Router = mux.NewRouter()

	// setup routes

	// api
	s.Router.HandleFunc("/ping", s.handlePing())
	s.Router.HandleFunc("/api/search", s.restEndpoint(s.handleSearch(false)))
	s.Router.HandleFunc("/api/product/{sku}", s.restEndpoint(s.handleProduct(false)))

	// app
	s.Router.HandleFunc("/app/search", s.handleSearch(true))
	s.Router.HandleFunc("/app/product/{sku}", s.handleProduct(true))

	s.HttpServer = &http.Server{
		Handler:      s.Router,
		Addr:         fmt.Sprintf("127.0.0.1:%d", s.Config.APIPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// initialize channels
	s.EventChan = make(chan Event, 0)

	return &s
}

func (s *Server) Start() {

	// connect to Couchbase buckets
	//cluster := MustOpenCluster("couchbase://localhost", "admin", "password123")
	cluster := MustOpenCluster(s.Config.CouchbaseClusterAddress, s.Config.CouchbaseUserName, s.Config.CouchbasePassword)
	pbucket := MustOpenBucket("products", cluster)
	tbucket := MustOpenBucket("tracking", cluster)
	sbucket := MustOpenBucket("sessions", cluster)

	s.CouchbaseCluster = cluster
	s.ProductsBucket = pbucket
	s.TrackingBucket = tbucket
	s.SessionsBucket = sbucket

	go s.processEvents()

	s.HttpServer.ListenAndServe()

}

func (s *Server) Shutdown() error {

	// sleep briefly to give processEvents() a chance to drain
	time.Sleep(time.Second * 1)

	err := s.TrackingBucket.Close()
	if err != nil {
		return err
	}
	err = s.ProductsBucket.Close()
	if err != nil {
		return err
	}
	err = s.CouchbaseCluster.Close()
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.HttpServer.Shutdown(ctx)
	log.Println("Shutdown server gracefully")
	return nil
}

func (s *Server) ServerReady() bool {

	var client = &http.Client{
		Timeout: time.Second * 5,
	}

	tries := 0
	for {
		tries++
		url := fmt.Sprintf("http://localhost:%d/ping", s.Port)
		resp, _ := client.Get(url)
		if resp != nil && resp.StatusCode == http.StatusOK {
			return true
		} else {
			time.Sleep(time.Second * 1)
		}
		if tries >= 5 {
			panic("Unable to connect to test HTTP server at " + url)
		}

	}

}
