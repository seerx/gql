package object

import (
	"errors"

	"github.com/graphql-go/graphql"
)

// MutationType 操作类型
var MutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addGoods": &graphql.Field{
				Type: goodsType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: goodsInputType,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input, isOK := p.Args["input"].(map[string]interface{})
					if !isOK {
						err := errors.New("Field 'addGoods' is missing required arguments: input. ")
						return nil, err
					}
					result := Goods{
						ID:    "random",
						Name:  input["name"].(string),
						Price: input["price"].(float64),
						URL:   input["url"].(string),
					}
					// 处理数据
					return result, nil
				},
			},
		},
	},
)
