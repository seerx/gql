package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"

	"github.com/seerx/gql"
	// _ "github.com/seerx/gql/examples"

	"github.com/graphql-go/handler"
)

func runGraphQL() {
	g := gql.Get()

	// gql.Get().RegisterInject(examples.Sser)

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

func main() {
	// testAST()
	runGraphQL()
}

func testAST() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset,
		"./object/json.go",
		nil,
		parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, d := range f.Decls {
		s, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		if s.Tok != token.TYPE {
			continue
		}
		// fmt.Println(s.Tok)
		for _, o := range s.Specs {
			_, ok := o.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// fmt.Println(t.Name)
		}

	}

	// ast.Print(fset, f)
}
