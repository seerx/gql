package gql

import (
	"errors"
	"fmt"
	"reflect"
)

type typeProp struct {
	TypeName    string
	PackageName string
	Kind        reflect.Kind // 真正的类别
	Type        reflect.Type // 真正的类型
	structField reflect.StructField
	// KeyType     reflect.Type // 映射 key 类型

	IsPrimitive bool // 是否原生类型
	IsPtr       bool // 是否指针
	IsList      bool // 是否列表
}

func (p *typeProp) Key() string {
	return fmt.Sprintf("%s@%s", p.TypeName, p.PackageName)
}

func (p *typeProp) KeyList() string {
	return fmt.Sprintf("%s_list@%s", p.TypeName, p.PackageName)
}

// func (p *typeProp) isList() bool {
// 	return p.Kind == reflect.Slice
// }

// func (p *typeProp) isPtr() bool {
// 	return p.Kind == reflect.Ptr
// }

func (p *typeProp) String() string {
	return fmt.Sprintf("%s<%s>primitive:%t", p.TypeName, p.PackageName, p.IsPrimitive)
}

func parseProp(paramType reflect.Type) *typeProp {
	kind := paramType.Kind()
	isPtr := kind == reflect.Ptr
	isList := kind == reflect.Slice

	if isPtr || isList {
		// 指针、切片、数组
		out := paramType.Elem()
		p := parseProp(out)
		p.IsPtr = isPtr
		p.IsList = isList
		return p
	} else if kind == reflect.Map || kind == reflect.Array {
		// 映射
		panic(errors.New("Do not use map (or array) as params type"))
		// out := paramType.Elem()
		// return parseProp(out)
	}

	p := &typeProp{}

	p.TypeName = paramType.Name()
	p.PackageName = paramType.PkgPath()
	p.IsPrimitive = p.PackageName == ""
	p.Type = paramType
	p.Kind = kind

	return p
}

type param struct {
	Name     string
	JSONName string
	Prop     *typeProp
	// IsResolveParams bool // 是否 graphql.ResolveParams 类型
	// KeyProp *typeProp // 映射时可用
	// IsPtr   bool      // 是否指针类型
	// IsList  bool      // 是否切片
}

func (p *param) String() string {
	tp := "simple"
	if p.Prop.IsList {
		tp = "list"
	} else if p.Prop.IsPtr {
		tp = "ptr"
	}

	return fmt.Sprintf("%s->%s\n%s", p.Name, tp, p.Prop)
}
