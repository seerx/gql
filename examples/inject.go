package examples

import (
	"context"
	"net/http"

	"github.com/seerx/gql/gql"
)

func init() {
	gql.Get().RegisterInject(Sser)
}

// Ss 测试
type Ss struct {
	A string
}

// Sser 测试
func Sser(ctx context.Context, r *http.Request) *Ss {
	// if ctx != nil {
	// 	panic(errors.New("Need Login"))
	// }
	return &Ss{
		A: "123456890",
	}
}
