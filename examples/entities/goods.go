package entities

import "time"

// Goods 商品信息
type Goods struct {
	ID       string    `json:"id" gql:"desc=商品 ID,lt=10"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	URL      string    `json:"url"`
	Time     time.Time `json:"time"`
	O        bool      `json:"o"`
	Ary      []string  `json:"ary"`
	Children []*Goods  `json:"children"`
}
