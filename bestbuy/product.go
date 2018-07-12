package bestbuy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"gopkg.in/couchbase/gocb.v1"
	gocbcore "gopkg.in/couchbase/gocbcore.v7"
)

type Product struct {
	Sku               string   `json:"sku"`
	ID                string   `json:"id"`
	Type              string   `json:"type"`
	Name              string   `json:"name"`
	RegPrice          float64  `json:"reg_price"`
	SalePrice         float64  `json:"sale_price"`
	Discount          float64  `json:"discount"`
	OnSale            string   `json:"on_sale"`
	ShortDescription  string   `json:"short_description"`
	Class             string   `json:"class"`
	BbItemID          string   `json:"bb_item_id"`
	ModelNumber       string   `json:"model_number"`
	Manufacturer      string   `json:"manufacturer"`
	Image             string   `json:"image"`
	MedImage          string   `json:"med_image"`
	ThumbImage        string   `json:"thumb_image"`
	LargeImage        string   `json:"large_image"`
	LongDescription   string   `json:"long_description"`
	Keywords          string   `json:"keywords"`
	CatDescendentPath string   `json:"cat_descendent_path"`
	CatIds            []string `json:"cat_ids"`
	Platform          string   `json:"platform"`
}

func GetProducts(prodIds []string, prodBucket *gocb.Bucket) []*Product {
	products := make([]*Product, len(prodIds))
	for i, key := range prodIds {
		//var prod Product
		prod, err := GetProduct(key, prodBucket) //prodBucket.Get(key, &prod)
		if _, ok := err.(*gocbcore.KvError); ok {
			prod.Name = "Product Not Found"
		}
		products[i] = prod
	}
	return products
}

func GetProductsFromHits(hits []gocb.SearchResultHit, prodBucket *gocb.Bucket) []*Product {
	products := make([]*Product, len(hits))
	for i, key := range hits {
		prod, err := GetProduct(key.Id, prodBucket)
		if _, ok := err.(*gocbcore.KvError); ok {
			prod.Name = "Product Not Found"
		}
		products[i] = prod
	}
	return products
}

func LoadProductsFromFile(pathToJson string) error {

	// open the bucket
	cluster, err := gocb.Connect("couchbase://localhost")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "admin",
		Password: "password123",
	})
	if err != nil {
		return err
	}
	bucket, err := cluster.OpenBucket("products", "")
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(pathToJson)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(file)
	for {
		s, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		product := Product{}
		err = json.Unmarshal([]byte(s), &product)
		if err != nil {
			return err
		}
		err = SaveProduct(product, bucket)
		if err != nil {
			return err
		}
	}
	err = bucket.Close()
	return err
}

func SaveProduct(product Product, bucket *gocb.Bucket) error {
	_, err := bucket.Upsert(product.ID, product, 0)
	return err
}

func GetProduct(ID string, bucket *gocb.Bucket) (*Product, error) {
	var product Product
	_, err := bucket.Get(ID, &product)
	return &product, err
}
