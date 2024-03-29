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
	// entities.Goods
	// 注册
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
