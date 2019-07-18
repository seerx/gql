package main

import (
	"github.com/seerx/gql/examples/util"
)

func init()  {
	registerGQL()
}

func main() {
	util.Start(8080)
}
