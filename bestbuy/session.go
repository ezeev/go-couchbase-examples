package bestbuy

type Session struct {
	ID            string         `json:"id"`
	SessionID     string         `json:"sessionId"`
	Queries       map[string]int `json:"queries"`
	ProductViews  map[string]int `json:"productViews"`
	CategoryViews map[string]int `json: "categoryViews"`
	BrandViews    map[string]int `json:"brandViews"`
	PlatformViews map[string]int `json:"platformViews"`
}

func NewSession(iD string) Session {
	s := Session{
		ID:            iD,
		SessionID:     iD,
		Queries:       make(map[string]int),
		ProductViews:  make(map[string]int),
		CategoryViews: make(map[string]int),
		BrandViews:    make(map[string]int),
		PlatformViews: make(map[string]int),
	}
	return s
}

/*type UserQuery struct {
	Timestamp int64  `json:"timestamp"`
	Query     string `json:"query"`
}

type UserProductView struct {
	Timestamp int64  `json:"timestamp"`
	ProductID string `json:"productID"`
}*/
