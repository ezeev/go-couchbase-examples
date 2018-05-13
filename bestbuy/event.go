package bestbuy

import (
	"fmt"
	"strings"
)

type Event interface {
	ID() string
	EventType() string
}

type QueryTrackEvent struct {
	id    string
	Type  string `json:"type"`
	Query string `json:"query"`
	Hour  int64  `json:"hour"`
	Day   int64  `json:"day"`
	Count int64  `json:"count"`
}

func (q QueryTrackEvent) ID() string {
	return q.id
}

func (q QueryTrackEvent) EventType() string {
	return q.Type
}

func NewQueryTrackEvent(query string) QueryTrackEvent {

	id := strings.ToLower(query)
	id = strings.Replace(id, " ", "_", -1)
	id = strings.Replace(id, "\"", "", -1)
	day, hour := EpochDayHour()

	id = fmt.Sprintf("%d-%d-%s", day, hour, id)

	//clean quotes off of query as well
	query = strings.Replace(query, "\"", "", -1)

	event := QueryTrackEvent{
		id:    id,
		Type:  "queryTrack",
		Query: query,
		Day:   day,
		Hour:  hour,
		Count: 1,
	}
	return event

}
