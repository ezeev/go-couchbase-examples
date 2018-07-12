package bestbuy

type Config struct {
	APIPort                 int    `json: "apiPort"`
	AssetsBaseUrl           string `json: "assetsBaseUrl"`
	ProductsFTSIndexName    string `json: "productsFTSIndexName"`
	TrackingFTSIndexName    string `json: "trackingFTSIndexName"`
	CouchbaseClusterAddress string `json: "couchbaseClusterAddress"`
	CouchbaseUserName       string `json: "couchbaseUserName"`
	CouchbasePassword       string `json: "couchbasePassword"`
}
