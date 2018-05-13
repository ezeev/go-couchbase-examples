package bestbuy

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) restEndpoint(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//do stuff
		q := r.URL.Query().Get("q")

		//execute a query
		res, err := Search("bb-search", q, s.ProductsBucket)
		if err != nil {
			log.Printf("Error: %s", err)
		}
		json.NewEncoder(w).Encode(res)
		//track the search
		event := NewQueryTrackEvent(q)
		s.EventChan <- event

	}
}
