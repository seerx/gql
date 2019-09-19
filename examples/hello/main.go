package main

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/seerx/gql"
)

// Hello 示例
type Hello struct {
	Message string
}

func Test1() (string, error) {
	return "", nil
}

func (Hello) Test1() (string, error) {
	return "", nil
}

func init() {
	gql.Get().RegisterQuery(func() (*Hello, error) {
		return &Hello{
			Message: "Hello GQL!",
		}, nil
	})
	gql.Get().RegisterMutation(Test1)
	gql.Get().RegisterQuery(Hello{})
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
