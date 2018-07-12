package bestbuy

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"

	"gopkg.in/couchbase/gocb.v1/cbft"

	"github.com/gorilla/mux"
	gocbcore "gopkg.in/couchbase/gocbcore.v7"
)

func (s *Server) restEndpoint(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	}
}

func (s *Server) handleSearch(web bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := r.URL.Query()
		q := params.Get("q")
		matchQuery := cbft.NewMatchQuery(q)
		matchQuery.Analyzer("en") // must specify the analyzer if you want stemming to work
		disQuery := cbft.NewDisjunctionQuery(matchQuery)

		autoboost := params.Get("autoboost")

		if autoboost != "" {
			//signals, err := SignalSearch("clicksfts", q, 100, s.TrackingBucket)
			signals, err := SignalSearch(s.Config.TrackingFTSIndexName, q, 100, s.TrackingBucket)
			if err != nil {
				log.Printf("Error getting signals: %s", err.Error())
			}
			if signals != nil {
				for k, v := range signals {
					boost := (math.Log(float64(v)) / 10) + 0.01
					bq := cbft.NewMatchQuery(k)
					bq.Boost(float32(boost))
					bq.Field("sku")
					disQuery.Or(bq)
				}
			}
		}

		// can we filters it
		/*boolQuery := cbft.NewBooleanQuery()
		fq := cbft.NewMatchQuery("Apple®")
		fq.Field("manufacturer")
		boolQuery.Must(fq)
		mainQuery.Or(boolQuery)*/

		sessionID := params.Get("session")
		var session Session
		//load session?
		if sessionID != "" {

			_, err := s.SessionsBucket.Get(sessionID, &session)
			if _, ok := err.(*gocbcore.KvError); ok {
				// new session
				session = NewSession(sessionID)
			} else {
				// existing session - personalize
				// brand
				for k, v := range session.BrandViews {
					bq := cbft.NewMatchQuery(k)
					bq.Field("manufacturer")
					bq.Boost(float32(v) / 100)
					disQuery.Or(bq)
				}
				// platform
				for k, v := range session.PlatformViews {
					bq := cbft.NewMatchQuery(k)
					bq.Field("platform")
					bq.Boost(float32(v) / 100)
					disQuery.Or(bq)
				}
			}
		}

		//execute a query
		//res, err := Search("productsfts", disQuery, s.ProductsBucket)
		res, err := Search(s.Config.ProductsFTSIndexName, disQuery, s.ProductsBucket)
		if err != nil {
			log.Printf("Error: %s", err)
		}
		res.Q = q
		res.Session = sessionID
		//res.Track =

		if !web {
			json.NewEncoder(w).Encode(res)
		} else {
			tmpl, err := template.ParseFiles("../templates/search.html", "../templates/header.html", "../templates/footer.html")
			if err != nil {
				panic(err)
			}
			err = tmpl.Execute(w, res)
			if err != nil {
				panic(err)
			}
		}
		//track the search
		event := NewQueryTrackEvent(q, res.TotalHits)
		s.EventChan <- event

		// update session?
		if sessionID != "" && q != "" {
			session.Queries[q] += 1
			_, err = s.SessionsBucket.Upsert(sessionID, session, 0)
			if err != nil {
				log.Printf("Error: %s", err)
			}
		}

	}
}

func (s *Server) handleProduct(web bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if sku, ok := vars["sku"]; ok {
			log.Printf("request for product: %s", sku)
			var prod Product
			_, err := s.ProductsBucket.Get(sku, &prod)
			if err != nil {
				log.Printf("Error: %s", err.Error())
			} else {
				if !web {
					json.NewEncoder(w).Encode(prod)
				} else {
					tmpl, err := template.ParseFiles("../templates/product.html", "../templates/header.html", "../templates/footer.html")
					if err != nil {
						panic(err)
					}
					err = tmpl.Execute(w, prod)
					if err != nil {
						panic(err)
					}
				}
			}

			//record event?
			params := r.URL.Query()
			if params.Get("track") == "true" {
				// build event
				query := params.Get("q")
				var catId string
				if len(prod.CatIds) > 0 {
					catId = prod.CatIds[len(prod.CatIds)-1]
				}
				event := NewClickTrackEvent(query, sku, catId)
				s.EventChan <- event

				// record session event?
				if params.Get("session") != "" {
					sessionID := params.Get("session")
					// does the session already exist?
					var session Session
					_, err := s.SessionsBucket.Get(sessionID, &session)
					if _, ok := err.(*gocbcore.KvError); ok {
						// new session
						session = NewSession(sessionID)
					}
					if len(prod.CatIds) > 0 {
						session.CategoryViews[prod.CatIds[len(prod.CatIds)-1]] += 1
					}
					session.PlatformViews[prod.Platform] += 1
					session.BrandViews[prod.Manufacturer] += 1
					session.ProductViews[prod.Sku] += 1

					// save
					_, err = s.SessionsBucket.Upsert(sessionID, session, 0)
					if err != nil {
						log.Printf("Error: %s", err)
					}

				}

			}

		} else {
			fmt.Fprintf(w, "SKU Id not provided")
		}

	}
}
