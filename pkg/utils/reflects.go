package utils

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// TypeProp 类型属性
type TypeProp struct {
	SrcType     reflect.Type // 原始类型
	TypeName    string       // 真正的类型名称
	PackageName string       // 包名称
	Kind        reflect.Kind // 真正的类型
	RealType    reflect.Type // 真正的类型，对于 切片和指针 来说
	// structField reflect.StructField
	// KeyType     reflect.Type // 映射 key 类型

	IsPrimitive bool // 是否原生类型
	IsPtr       bool // 是否指针
	IsList      bool // 是否列表，切片
}

// Key 类型唯一标识
func (t *TypeProp) Key() string {
	return fmt.Sprintf("%s@%s", t.TypeName, t.PackageName)
}

// KeyOfList 切片类型唯一标识
func (t *TypeProp) KeyOfList() string {
	return fmt.Sprintf("%s_list@%s", t.TypeName, t.PackageName)
}

// String 转为字符串
func (t *TypeProp) String() string {
	return fmt.Sprintf("%s<%s>primitive:%t", t.TypeName, t.PackageName, t.IsPrimitive)
}

// IsContext 是否是 context.Context 接口
func (t *TypeProp) IsContext() bool {
	return t.PackageName == "context" && t.TypeName == "Context"
}

// ParseTypeProp 解析类型属性
func ParseTypeProp(typ reflect.Type) *TypeProp {
	var p = &TypeProp{SrcType: typ, RealType: typ, Kind: typ.Kind()}
	isPtr := p.Kind == reflect.Ptr
	isList := p.Kind == reflect.Slice
	pkg := typ.PkgPath()

	p.PackageName = pkg
	p.TypeName = typ.Name()

	if isPtr || isList {
		// 指针或切片
		p = ParseTypeProp(typ.Elem())
		p.SrcType = typ
		p.IsList = isList
		p.IsPtr = isPtr
		// p.PackageName = pkg
	} else if p.Kind == reflect.Map || p.Kind == reflect.Array {
		// 映射或数组，不支持
		panic(errors.New("Do not use map (or array) as params type"))
	}

	p.IsPrimitive = p.PackageName == ""

	return p
}

// FuncInfo 函数信息
type FuncInfo struct {
	Name string // 名称
	Pkg  string // 包
}

// func getFuncName(aFunc interface{}, seps ...rune) string {
// 	// 获取函数名称
// 	fn := runtime.FuncForPC(reflect.ValueOf(aFunc).Pointer()).Name()

// 	// 用 seps 进行分割
// 	fields := strings.FieldsFunc(fn, func(sep rune) bool {
// 		for _, s := range seps {
// 			if sep == s {
// 				return true
// 			}
// 		}
// 		return false
// 	})

// 	if size := len(fields); size > 0 {
// 		return fields[size-1]
// 	}
// 	return ""
// }

// func getFuncPackage(aFunc interface{}) string {
// 	// 获取函数名称
// 	fn := runtime.FuncForPC(reflect.ValueOf(aFunc).Pointer()).Name()

// 	p := strings.LastIndex(fn, ".")

// 	if p > 0 {
// 		return fn[0:p]
// 	}

// 	return ""
// }

// ParseFuncInfo 解析函数信息
func ParseFuncInfo(aFunc interface{}) *FuncInfo {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(aFunc).Pointer()).Name()

	p := strings.LastIndex(fn, ".")

	if p > 0 {

		return &FuncInfo{
			Name: fn[p+1:],
			Pkg:  fn[0:p],
		}
	}

	return nil
}

//FullName 函数的完整名称
func (f *FuncInfo) FullName() string {
	return f.Pkg + "." + f.Name
}

// ParseStructFieldName 解析结构字段名称
func ParseStructFieldName(field *reflect.StructField) string {
	name := field.Tag.Get("json")
	if name == "" {
		if field.Anonymous {
			// 匿名字段
			name = field.Name
		} else {
			name = field.Name
		}
		// tmp := strings.ToLower(string(name[0]))
		// name = strings.Replace(name, string(name[0]), tmp, 1)
	}
	dotp := strings.Index(name, ",")
	if dotp >= 0 {
		return name[0:dotp]
	}
	return name
}
