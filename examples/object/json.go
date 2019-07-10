package object

import (
	"time"

	"github.com/graphql-go/graphql"
)

// Goods 商品信息
type Goods struct {
	ID    string    `json:"id" gql:"aaa"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	URL   string    `json:"url"`
	Time  time.Time `json:"time"`
	O     bool      `json:"o"`
}

var goodsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Goods",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var goodsListType = graphql.NewList(goodsType)

var goodsInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "goodsInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"price": &graphql.InputObjectFieldConfig{
				Type: graphql.Float,
			},
			"url": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)
