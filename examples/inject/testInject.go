package main

import (
	"fmt"

	"github.com/seerx/gql"
	"github.com/seerx/gql/examples/entities"
)

func init() {
	gql.Get().RegisterQuery(UseInject)
}

// UseInject 测试注入信息
func UseInject(a *Account, ai A) (*entities.Goods, error) {
	fmt.Printf("Inject value: %s\n", a)
	ai.Do("nonono")
	return &entities.Goods{
		ID:    a.String(),
		Name:  "热水器",
		Price: 30.8,
		URL:   "http://www.sohu.com",
	}, nil
}
