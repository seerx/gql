package util

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/seerx/gql"
)

// Start Start server
func Start(port int) {
	g := gql.Get()
	handler := g.NewHandler(&handler.Config{
		Pretty:   true,
		GraphiQL: true,
	})

	fmt.Print(g.Summary())

	http.Handle("/graphql", handler)
	fmt.Println("The api server will run on port : ", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
