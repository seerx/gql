package main

import (
	"github.com/seerx/gql"
	"github.com/seerx/gql/examples/entities"
)

/**
mutation
{
    addGoods(input:{name:"123", price:10.9,url:"www.sohu.com"}) {
      id, name,url,price
    }
}
*/
func init() {
	g := gql.Get()
	// 随意注册一个查询操作，否则会报错
	g.RegisterQuery(func() (entities.Hello, error) {
		return entities.Hello{
			Message: "Hello GQL",
		}, nil
	})
	// 注册 TestMutation 方法为 mutation 操作
	g.RegisterMutation(TestMutation)
}

// TestMutation 测试操作
func TestMutation(param *entities.Goods) (*entities.Goods, error) {
	result := entities.Goods{
		ID:    "random",
		Name:  param.Name,
		Price: param.Price,
		URL:   param.URL,
	}
	// 处理数据
	return &result, nil
}
