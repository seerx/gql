package main

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/seerx/gql"
	"github.com/seerx/gql/examples/entities"
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

// GoodsQuery 货物查询
type GoodsQuery struct{}

func init() {
	g := gql.Get()
	// 注册
	// 把 GoodsQuery 结构的所有方法注册为 Query 操作
	g.RegisterQuery(GoodsQuery{})
	// 把 AGoodsList 方法注册为 Query 操作
	g.RegisterQuery(AGoodsList)
}

// GoodsList 获取货物列表
func (good GoodsQuery) GoodsList() ([]entities.Goods, error) {
	return []entities.Goods{
		entities.Goods{
			ID: "1", Name: "test1", Time: time.Now(), O: false,
		},
		entities.Goods{
			ID: "2", Name: "test2", Time: time.Now(), O: true,
		},
	}, nil
}

// Goods 查询
func (good GoodsQuery) Goods(param *entities.Goods) (*entities.Goods, error) {
	// fmt.Printf("Inject value: %s\n", s.A)
	return &entities.Goods{
		ID:    param.ID,
		Name:  "热水器",
		Price: 30.8,
		URL:   "http://www.sohu.com",
	}, nil
}

// Test sd
// func (good GoodsQuery) Test(v *gqlh.InputValidator, in *entities.Goods) (*entities.Goods, error) {
// 	v.Requires("id")

// 	return &entities.Goods{
// 		ID:    in.ID,
// 		Name:  "热水器",
// 		Price: 30.8,
// 		URL:   "http://www.sohu.com",
// 	}, nil
// }

// AGoodsList 单个函数
func AGoodsList(p graphql.ResolveParams) ([]entities.Goods, error) {
	return []entities.Goods{
		entities.Goods{
			ID: "A1", Name: "A-test1",
		},
		entities.Goods{
			ID: "A2", Name: "A-test2",
		},
	}, nil
}
