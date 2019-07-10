package main

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/seerx/gql"
	"github.com/seerx/gql/examples/entities"
)

// Hello 示例
type Hello struct {
	Message string
}

func init() {
	gql.Get().RegisterQuery(func() (*entities.Hello, error) {
		return &entities.Hello{
			Message: "Hello GQL!",
		}, nil
	})
}

func main() {
	g := gql.Get()
	handler := g.NewHandler(&handler.Config{
		Pretty:   true,
		GraphiQL: true,
	})

	fmt.Print(g.Summary())

	apiPort := 8080

	http.Handle("/graphql", handler)
	fmt.Println("The api server will run on port : ", apiPort)
	http.ListenAndServe(fmt.Sprintf(":%d", apiPort), nil)
}
