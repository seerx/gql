package gqlh

import (
	"io"
	"reflect"

	"github.com/graphql-go/graphql"
)

//函数输入参数定义
type inputParam interface {
	createValue(structInstance interface{},
		gqlPara *graphql.ResolveParams,
		validator *InputValidator,
		requestValue reflect.Value) reflect.Value
	isInjectInterface() bool
}

// 结构实例参数，即函数属于结构体时，其第一个参数
type ipStruct struct {
}

func (*ipStruct) createValue(structInstance interface{},
	gqlPara *graphql.ResolveParams,
	validator *InputValidator,
	requestValue reflect.Value) reflect.Value {
	return reflect.ValueOf(structInstance)
}

func (*ipStruct) isInjectInterface() bool {
	return false
}

// Graphql 原始参数
type ipGraphqlResolveParams struct {
	valueIsPtr bool
}

func (ipgrp *ipGraphqlResolveParams) createValue(structInstance interface{},
	gqlParam *graphql.ResolveParams,
	validator *InputValidator,
	requestValue reflect.Value) reflect.Value {
	if ipgrp.valueIsPtr {
		return reflect.ValueOf(gqlParam)
	}
	return reflect.ValueOf(*gqlParam)
}

func (ipgrp *ipGraphqlResolveParams) isInjectInterface() bool {
	return false
}

// 参数验证参数
type ipValidator struct {
}

func (*ipValidator) createValue(structInstance interface{},
	gqlParam *graphql.ResolveParams,
	validator *InputValidator,
	requestValue reflect.Value) reflect.Value {
	return reflect.ValueOf(validator)
}

func (*ipValidator) isInjectInterface() bool {
	return false
}

// 注入参数
type ipInject struct {
	valueIsInterface bool
	inject           *InjectInfo
}

func (ipi *ipInject) createValue(structInstance interface{},
	gqlParam *graphql.ResolveParams,
	validator *InputValidator,
	requestValue reflect.Value) reflect.Value {
	return ipi.inject.CallFn(gqlParam)
}
func (ipi *ipInject) isInjectInterface() bool {
	return ipi.valueIsInterface
}

// 提交参数
type ipRequest struct {
}

func (*ipRequest) createValue(structInstance interface{},
	gqlParam *graphql.ResolveParams,
	validator *InputValidator,
	requestValue reflect.Value) reflect.Value {
	return requestValue
}

func (*ipRequest) isInjectInterface() bool {
	return false
}

// IsCloser 判断该值是否时实现了 io.Closer 接口
// 如果实现，则返回 ioCloser 结构
// 否则返回 nil
func IsCloser(value reflect.Value) io.Closer {
	val := value.Interface()
	if val != nil {
		closer, ok := val.(io.Closer)
		if ok {
			return closer
		}
	}
	return nil
}
