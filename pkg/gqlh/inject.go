package gqlh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/seerx/gql/pkg/utils"
)

// Inject 注入结构
type Inject struct {
	injectMap map[reflect.Type]*InjectInfo
}

// NewInject 创建注入管理器
func NewInject() *Inject {
	return &Inject{
		injectMap: make(map[reflect.Type]*InjectInfo),
	}
}

const keyOfConext = "context"
const keyOfRequest = "request"

var errOfInject = errors.New("Param injectFn must be a func  like -- func(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *returnType, and returnType must be a struct or interface")
var typeOfHTTPRequest = reflect.TypeOf(http.Request{})
var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
var typeOfGqlParams = reflect.TypeOf(graphql.ResolveParams{})

//var typeOfHTTPRequest = reflect.TypeOf(http.Request{})

// InjectInfo 注入对象结构
type InjectInfo struct {
	ReturnType reflect.Type    // 返回类型
	Fn         reflect.Value   // 生成注入对象的函数
	FnInfo     *utils.FuncInfo // 函数的信息
}

// CallFn 调用注入函数
func (i *InjectInfo) CallFn(p *graphql.ResolveParams) reflect.Value {
	var ok bool
	root, ok := p.Info.RootValue.(map[string]interface{})
	if !ok {
		panic(errors.New("No root value found"))
	}

	ctx, ok := root[keyOfConext]
	if !ok {
		panic(errors.New("No context.Context int root value"))
	}
	req, ok := root[keyOfRequest]
	if !ok {
		panic(errors.New("No http.Request int root value"))
	}

	args := make([]reflect.Value, 3)
	args[0] = reflect.ValueOf(ctx)
	args[1] = reflect.ValueOf(req)
	args[2] = reflect.ValueOf(p)

	// 调用注入函数
	res := i.Fn.Call(args)

	return res[0]
}

// FindInject 根据类型查找注入信息
func (ij *Inject) FindInject(typ reflect.Type) *InjectInfo {
	info, ok := ij.injectMap[typ]
	if ok {
		return info
	}
	return nil
}

// StoreContext 存储Conext 和 Request，以备注入时使用
func (ij *Inject) StoreContext(ctx context.Context, r *http.Request, root map[string]interface{}) {
	root[keyOfConext] = ctx
	root[keyOfRequest] = r
}

func checkInputArg(arg reflect.Type, target reflect.Type) {
	inp := utils.ParseTypeProp(arg)
	if inp.Kind == reflect.Interface {
		// 接口类型
		tp := utils.ParseTypeProp(target)
		if inp.RealType != tp.RealType {
			// 类型不对
			panic(errOfInject)
		}
	} else {
		if !inp.IsPtr { // 不是指针
			panic(errOfInject)
		}
		if inp.RealType != target {
			// 类型不对
			panic(errOfInject)
		}
	}
}

// Inject 注入函数
func (ij *Inject) Inject(injectFn interface{}) {
	typ := reflect.TypeOf(injectFn)
	if typ.Kind() != reflect.Func {
		panic(errOfInject)
	}

	// 取得返回值类型
	outCount := typ.NumOut()
	if outCount != 1 { // 返回参数必须只能有一个
		panic(errOfInject)
	}
	out := typ.Out(0)
	outp := utils.ParseTypeProp(out)

	if outp.Kind != reflect.Interface {
		if !outp.IsPtr { // 返回值不是指针类型
			panic(errOfInject)
		}
		if outp.IsPrimitive { // 返回值是原始类型
			panic(errOfInject)
		}
		if outp.Kind != reflect.Struct { // 返回值不是结构类型
			panic(errOfInject)
		}
	}

	inject, ok := ij.injectMap[outp.RealType]
	if ok { // 已经存在
		panic(fmt.Errorf("Type [%s] is Registered by func [%s]", outp.TypeName, inject.FnInfo.FullName()))
	}

	// 判断输入参数
	inCount := typ.NumIn()
	if inCount != 3 {
		panic(errOfInject)
	}
	// 第 1 个参数 context.Context
	//in := typ.In(0)
	checkInputArg(typ.In(0), typeOfContext)

	//inp := utils.ParseTypeProp(in)
	//if !inp.IsContext() {
	//	// 类型不对
	//	panic(errOfInject)
	//}
	// 第 2 个参数 *http.Request
	//in = typ.In(1)
	checkInputArg(typ.In(1), typeOfHTTPRequest)
	//inp = utils.ParseTypeProp(in)
	//if !inp.IsPtr { // 不是指针
	//	panic(errOfInject)
	//}
	//if inp.RealType != typeOfHTTPRequest {
	//	// 类型不对
	//	panic(errOfInject)
	//}

	// 第 3 个参数 *graphql.ResolveParams
	checkInputArg(typ.In(2), typeOfGqlParams)

	// 注册函数
	fnInfo := utils.ParseFuncInfo(injectFn)
	inject = &InjectInfo{
		ReturnType: outp.RealType,
		FnInfo:     fnInfo,
		Fn:         reflect.ValueOf(injectFn),
	}
	ij.injectMap[inject.ReturnType] = inject
}
