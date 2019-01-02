package model

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
