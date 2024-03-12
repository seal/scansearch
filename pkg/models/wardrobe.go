package models

import "gorm.io/gorm"

type Wardrobe struct {
	gorm.Model
	ID                int     `gorm:"primary_key,omitempty"`
	UserID            int     `json:"userid,omitempty"`
	Like              bool    `json:"like,omitempty"` // true = like, false = love
	DesiredPrice      float64 `json:"desiredprice,omitempty"`
	SerpapiProductApi string  `json:"serpapi_product_api,omitempty"`
	ImageURL          string  `json:"image_url,omitempty"`

	Position       int     `json:"position,omitempty"`
	Title          string  `json:"title,omitempty"`
	Link           string  `json:"link,omitempty"`
	ProductLink    string  `json:"product_link,omitempty"`
	ProductID      string  `json:"product_id,omitempty"`
	Source         string  `json:"source,omitempty"`
	Price          string  `json:"price,omitempty"`
	ExtractedPrice float64 `json:"extracted_price,omitempty"`
	Rating         float64 `json:"rating,omitempty"`
	Reviews        int     `json:"reviews,omitempty"`
	//Extensions     interface{} `json:"extensions,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Delivery  string `json:"delivery,omitempty"`
}
type WardrobeResponse struct {
	ID                int     `json:"wardrobid,omitempty"`
	UserID            int     `json:"userid,omitempty"`
	Like              bool    `json:"like,omitempty"` // true = like, false = love
	DesiredPrice      float64 `json:"desiredprice,omitempty"`
	SerpapiProductApi string  `json:"serpapi_product_api,omitempty"`
	ImageURL          string  `json:"image_url,omitempty"`

	Position       int     `json:"position,omitempty"`
	Title          string  `json:"title,omitempty"`
	Link           string  `json:"link,omitempty"`
	ProductLink    string  `json:"product_link,omitempty"`
	ProductID      string  `json:"product_id,omitempty"`
	Source         string  `json:"source,omitempty"`
	Price          string  `json:"price,omitempty"`
	ExtractedPrice float64 `json:"extracted_price,omitempty"`
	Rating         float64 `json:"rating,omitempty"`
	Reviews        int     `json:"reviews,omitempty"`
	//Extensions     interface{} `json:"extensions,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Delivery  string `json:"delivery,omitempty"`
	//ProductDetails    ProductDetails `json:"productDetails,omitempty"`
}
type ProductDetails struct {
	Position          int     `json:"position,omitempty"`
	Title             string  `json:"title,omitempty"`
	Link              string  `json:"link,omitempty"`
	ProductLink       string  `json:"product_link,omitempty"`
	ProductID         string  `json:"product_id,omitempty"`
	SerpapiProductAPI string  `json:"serpapi_product_api,omitempty"`
	Source            string  `json:"source,omitempty"`
	Price             string  `json:"price,omitempty"`
	ExtractedPrice    float64 `json:"extracted_price,omitempty"`
	Rating            float64 `json:"rating,omitempty"`
	Reviews           int     `json:"reviews,omitempty"`
	//Extensions        interface{} `json:"extensions,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Delivery  string `json:"delivery,omitempty"`
}
type WardrobeInput struct {
	Like              bool    `json:"like" binding:"required,omitempty"`
	DesiredPrice      float64 `json:"desiredprice,omitempty"`
	SerpapiProductApi string  `json:"serpapi_product_api,omitempty"`
	ImageURL          string  `json:"image_url,omitempty"`
	Position          int     `json:"position,omitempty"`
	Title             string  `json:"title,omitempty"`
	Link              string  `json:"link,omitempty"`
	ProductLink       string  `json:"product_link,omitempty"`
	ProductID         string  `json:"product_id,omitempty"`
	Source            string  `json:"source,omitempty"`
	Price             string  `json:"price,omitempty"`
	ExtractedPrice    float64 `json:"extracted_price,omitempty"`
	Rating            float64 `json:"rating,omitempty"`
	Reviews           int     `json:"reviews,omitempty"`
	//Extensions     interface{} `json:"extensions,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Delivery  string `json:"delivery,omitempty"`
	//ProductDetails    ProductDetails `json:"productDetails,omitempty"`
}

// Removed due to being the same as wardrobe input
/*
type WardrobePut struct {
	//ID                int     `json:"id" binding:"required,omitempty"`
	Like              bool    `json:"like,omitempty"`
	DesiredPrice      float64 `json:"desiredprice,omitempty"`
	SerpapiProductApi string  `json:"serpapi_product_api,omitempty"`
	ImageURL          string  `json:"image_url,omitempty"`
	Position          int     `json:"position,omitempty"`
	Title             string  `json:"title,omitempty"`
	Link              string  `json:"link,omitempty"`
	ProductLink       string  `json:"product_link,omitempty"`
	ProductID         string  `json:"product_id,omitempty"`
	Source            string  `json:"source,omitempty"`
	Price             string  `json:"price,omitempty"`
	ExtractedPrice    float64 `json:"extracted_price,omitempty"`
	Rating            float64 `json:"rating,omitempty"`
	Reviews           int     `json:"reviews,omitempty"`
	//Extensions     interface{} `json:"extensions,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Delivery  string `json:"delivery,omitempty"`
	//ProductDetails    ProductDetails `json:"productDetails,omitempty"`
}*/
