package bestbuy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"gopkg.in/couchbase/gocb.v1"
)

/*

{
	"sku_s": "18941592",
	"id": "18941592",
	"type_s": "Music",
	"name_s": "Moot and Lid - CD",
	"reg_price_f": 13.99,
	"sale_price_f": 13.99,
	"discount_f": 0.0,
	"on_sale_s": "false",
	"short_description_s": null,
	"short_description_t": null,
	"class_s": "COMPACT DISC",
	"class_t": "COMPACT DISC",
	"bb_item_id_s": "1548609",
	"model_number_s": "",
	"manufacturer_s": "Ictus Records",
	"manufacturer_t": "Ictus Records",
	"image_s": "http://images.bestbuy.com/BestBuy_US/images/products/1894/18941592.jpg",
	"med_image_s": null,
	"thumb_image_s": "http://images.bestbuy.com/BestBuy_US/images/products/1894/18941592s.jpg",
	"large_image_s": null, "long_description_t": "", "keywords_txt_en": "Ictus Records Moot and Lid - CD  None COMPACT DISC", "cat_descendent_path": "cat00000|Best Buy/abcat0600000|Movies & Music/cat02001|Music/cat02007|Jazz", "cat_id_ss": ["cat00000", "abcat0600000", "cat02001", "cat02007"]}

*/

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
}

func LoadProducts(pathToJson string) error {

	// open the bucket
	cluster, err := gocb.Connect("couchbase://localhost")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "evan",
		Password: "password",
	})
	if err != nil {
		return err
	}
	bucket, err := cluster.OpenBucket("bb-catalog", "")
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
	if err != nil {
		return err
	}
	return nil
}

func SaveProduct(product Product, bucket *gocb.Bucket) error {
	_, err := bucket.Upsert(product.ID, product, 0)
	if err != nil {
		return err
	}
	return nil
}
