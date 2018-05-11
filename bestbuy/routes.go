package bestbuy

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (s *Server) Start(port string) {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/api", s.handleHome())

	srv := &http.Server{
		Handler:      s.Router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//do stuff
		fmt.Fprint(w, "This is home.")
	}
}
