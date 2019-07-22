package utils

import (
	"fmt"
	"reflect"
	"time"

	"github.com/graphql-go/graphql"
)

// StructFieldTypeToGraphType 结构字段类型对应 graphql 类型
func StructFieldTypeToGraphType(field *reflect.StructField) (grapghType graphql.Output, isStruct bool) {
	tp := field.Type

	kind := tp.Kind()

	isPtr := kind == reflect.Ptr
	isList := kind == reflect.Slice

	if isPtr || isList {
		// 指针、切片、数组
		// field.Type.Elem()
		//kind = field.Type.Elem().Kind()
		tp = field.Type.Elem()
	}

	gType, isStruct, err := TypeToGraphQLType(tp)
	if err != nil {
		panic(err)
	}
	if isStruct {
		return nil, true
	}
	if isList {
		return graphql.NewList(gType), false
	}
	return gType, false
}

// TypeToGraphQLType go 类型转 graphql 类型
func TypeToGraphQLType(typ reflect.Type) (outType graphql.Output, isStruct bool, err error) {
	kind := typ.Kind()
	if IsIntType(kind) {
		return graphql.Int, false, nil
	}
	if IsFloatType(kind) {
		return graphql.Float, false, nil
	}
	if IsStringType(kind) {
		return graphql.String, false, nil
	}
	if IsBoolType(kind) {
		return graphql.Boolean, false, nil
	}
	if IsTimeType(typ) {
		return graphql.DateTime, false, nil
	}
	if IsStructType(kind) {
		return nil, true, nil
	}

	return nil, false, fmt.Errorf("不支持 %s 类型", typ.Name())
}

// IsIntType 是否整数类型
func IsIntType(kind reflect.Kind) bool {
	return kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64
}

// IsFloatType 是否浮点类型
func IsFloatType(kind reflect.Kind) bool {
	return kind == reflect.Float32 ||
		kind == reflect.Float64
}

// IsStringType 是否字符串类型
func IsStringType(kind reflect.Kind) bool {
	return kind == reflect.String
}

// IsBoolType 是否布尔类型
func IsBoolType(kind reflect.Kind) bool {
	return kind == reflect.Bool
}

// IsTimeType 是否时间类型
func IsTimeType(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct && typ == reflect.TypeOf(time.Time{})
}

// IsStructType 是否结构类型
func IsStructType(kind reflect.Kind) bool {
	return kind == reflect.Struct
}
