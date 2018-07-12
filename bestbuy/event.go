package bestbuy

import (
	"fmt"
	"log"
	"strings"

	gocbcore "gopkg.in/couchbase/gocbcore.v7"
)

type Event interface {
	ID() string
	EventType() string
}

// QueryTrackEvent is used to track queries from /api/search
type QueryTrackEvent struct {
	id      string
	Type    string `json:"type"`
	Query   string `json:"query"`
	Hour    int64  `json:"hour"`
	Day     int64  `json:"day"`
	Count   int64  `json:"count"`
	NumHits int    `json:"numHits"`
}

func (q QueryTrackEvent) ID() string {
	return q.id
}

func (q QueryTrackEvent) EventType() string {
	return q.Type
}

func NewQueryTrackEvent(query string, numHits int) QueryTrackEvent {

	id := strings.ToLower(query)
	id = strings.Replace(id, " ", "_", -1)
	id = strings.Replace(id, "\"", "", -1)
	day, hour := EpochDayHour()

	id = fmt.Sprintf("%d-%d-%s", day, hour, id)

	//clean quotes off of query as well
	query = strings.Replace(query, "\"", "", -1)

	event := QueryTrackEvent{
		id:      id,
		Type:    "queryTrack",
		Query:   query,
		Day:     day,
		Hour:    hour,
		Count:   1,
		NumHits: numHits,
	}
	return event

}

// ClickTrackEvent
type ClickTrackEvent struct {
	id       string
	Type     string `json:"type"`
	Query    string `json:"query"`
	Category string `json:"category"`
	Sku      string `json:"sku"`
	Hour     int64  `json:"hour"`
	Day      int64  `json:"day"`
	Count    int64  `json:"count"`
}

func (c ClickTrackEvent) ID() string {
	return c.id
}

func (c ClickTrackEvent) EventType() string {
	return c.Type
}

func NewClickTrackEvent(query string, sku string, category string) ClickTrackEvent {
	id := strings.ToLower(query)
	id = strings.Replace(id, " ", "_", -1)
	id = strings.Replace(id, "\"", "", -1)
	id += "-" + sku
	day, hour := EpochDayHour()

	id = fmt.Sprintf("%d-%d-%s", day, hour, id)

	e := ClickTrackEvent{
		id:       id,
		Type:     "clickTrack",
		Query:    query,
		Category: category,
		Sku:      sku,
		Hour:     hour,
		Day:      day,
		Count:    1,
	}
	return e
}

func (s *Server) processEvents() {
	for event := range s.EventChan {
		id := event.ID()
		//t := event.EventType()

		// Query Tracking
		//if t == "queryTrack" {
		_, err := s.TrackingBucket.Insert(id, event, 0)
		if err != nil {
			if _, ok := err.(*gocbcore.KvError); ok {
				//key already exists, increment counter
				s.TrackingBucket.MutateIn(id, 0, 0).Counter("count", 1, false).Execute()
				log.Printf("Incremented counter for id %s", id)
			}
		} else {
			log.Printf("Added counter for id %s", id)
		}
		//}
	}
}
