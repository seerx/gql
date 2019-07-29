package main

import (
	"fmt"

	"github.com/seerx/gql"
	"github.com/seerx/gql/examples/entities"
)

func init() {
	gql.Get().RegisterQuery(O{})
}

type O struct {
}

func (O) UseInjectDesc() string {
	return `测试注入功能`
}

// UseInject 测试注入信息
func (O) UseInject(a *Account, ai A) (*entities.Goods, error) {
	fmt.Printf("Inject value: %s\n", a)
	ai.Do("nonono")
	return &entities.Goods{
		ID:    a.String(),
		Name:  "热水器",
		Price: 30.8,
		URL:   "http://www.sohu.com",
		Children: []*entities.Goods{
			&entities.Goods{
				ID:    "child1",
				Name:  "OK1",
				Price: 22.0,
			},
			&entities.Goods{
				ID:    "child2",
				Name:  "OK3",
				Price: 22.4,
				Children: []*entities.Goods{
					&entities.Goods{
						ID:    "child2-1",
						Name:  "OK1111",
						Price: 22.0,
					},
					&entities.Goods{
						ID:    "child2-2",
						Name:  "OK3222",
						Price: 22.4,
					},
				},
			},
		},
	}, nil
}
