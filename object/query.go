package object

import (
	"errors"

	"github.com/graphql-go/graphql"
)

/*
{
  goodsList {
    id, name
  },
  goods(id:"100") {
    id,
    name,
    price
  }
}
**/

/**
一次调用多次执行同一个接口时，使用别名
{
  G1:Goods(Goods:{id:"123"}) {
    id, name, time, o
  }
  G2:Goods(Goods:{id:"321"}){
    id, name, time, o
  }
  GoodsList {
    id, time, name
  }
}*/

// QueryType 查询类型
var QueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"goodsList": &graphql.Field{
				Type: goodsListType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return []Goods{
						Goods{
							ID: "1", Name: "test1",
						},
						Goods{
							ID: "2", Name: "test2",
						},
					}, nil
				},
			},
			"goods": &graphql.Field{
				Type: goodsType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := p.Args["id"].(string)
					if isOK {
						return &Goods{
							ID:    idQuery,
							Name:  "热水器",
							Price: 30.8,
							URL:   "http://www.sohu.com",
						}, nil
					}

					err := errors.New("Field 'goods' is missing required arguments: id. ")
					return nil, err
				},
			},
		},
	},
)
