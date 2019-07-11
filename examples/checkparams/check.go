package main

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/seerx/gql"
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
	gql.Get().RegisterQueryWithValidateFn(require, func(paramName string, paramValue interface{}, inputParam graphql.ResolveParams) error {
		fmt.Println(paramName)
		// if "id" == paramName {
		// 	return errors.New("ID Test")
		// }
		return nil
	})
}

// Storage asa
type Storage struct {
	Addr string
	Code string
}

// goods 商品信息
type goods struct {
	ID    string    `json:"id" gql:"aaa"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	URL   string    `json:"url"`
	Time  time.Time `json:"time,omitempty"`
	O     bool      `json:"o"`
	S     *Storage
}

// require 必须的参数检查
func require(v *gqlh.InputValidator, in *goods) (*goods, error) {
	v.Requires("id")

	return &goods{
		ID:    in.ID,
		Name:  "热水器",
		Price: 30.8,
		URL:   "http://www.sohu.com",
		S: &Storage{
			Addr: in.S.Addr,
			Code: "00001",
		},
	}, nil
}
