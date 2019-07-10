package gqlh

import "github.com/seerx/gql/pkg/utils"

// Field 字段，参数，结构字段等等
type Field struct {
	Name     string
	JSONName string // 适用于结构字段
	Prop     *utils.TypeProp
}
