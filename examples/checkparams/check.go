package main

import (
	"github.com/seerx/gql"
	"github.com/seerx/gql/examples/entities"
	"github.com/seerx/gql/pkg/gqlh"
)

/*
以下请求会返回错误
{
  Require(Goods:{mame:"12"}) {
    id, name
  }
}

// 以下请求会完成执行
{
  Require(Goods:{id:"12"}) {
    id, name
  }
}
*/

func init() {
	gql.Get().RegisterQuery(Require)
}

// Require 必须的函数
func Require(v *gqlh.InputValidator, in *entities.Goods) (*entities.Goods, error) {
	v.Requires("id")

	return &entities.Goods{
		ID:    in.ID,
		Name:  "热水器",
		Price: 30.8,
		URL:   "http://www.sohu.com",
	}, nil
}
