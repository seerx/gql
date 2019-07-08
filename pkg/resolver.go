package pkg

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/seerx/gql/utils"
)

// Resolver 查询和操作函数定义
type Resolver struct {
	manager *ResolverManager
	// gql                    *GQL
	structInstance         interface{} // 结构内方法
	out                    *Field
	outError               *Field
	params                 []*Field
	executor               reflect.Value
	input                  *RequestObject // 输入参数
	isGraphQLParamInParams bool           // params 是否包含 GraphSQL resolve 参数
	isValidatorInParams    bool           // params 是否包含 InputValidator
}

// ResolverManager 管理器
type ResolverManager struct {
	name          string
	resObjManager *ResponseObjectManager
	reqObjManager *RequestObjectManager
	inject        *Inject
	resolverMap   map[string]*Resolver
}

// RegisterInfo 注册状态
type RegisterInfo struct {
	Type    string
	Package string
	Struct  string
	Func    string
	Error   string
}

// graphql 参数类型
var typeOfResolveParams = reflect.TypeOf(graphql.ResolveParams{})

// NewQueryResolverManager 创建查询操作管理器
func NewQueryResolverManager(inject *Inject,
	responseObjectManager *ResponseObjectManager,
	requestObjectManager *RequestObjectManager) *ResolverManager {
	return &ResolverManager{
		name:          "Query",
		inject:        inject,
		resObjManager: responseObjectManager,
		reqObjManager: requestObjectManager,
		resolverMap:   make(map[string]*Resolver),
	}
}

// NewMutationResolverManager 创建操作管理器
func NewMutationResolverManager(inject *Inject,
	responseObjectManager *ResponseObjectManager,
	requestObjectManager *RequestObjectManager) *ResolverManager {
	return &ResolverManager{
		name:          "Mutation",
		inject:        inject,
		resObjManager: responseObjectManager,
		reqObjManager: requestObjectManager,
		resolverMap:   make(map[string]*Resolver),
	}
}

// CreateResolveObject 创建查询对象结构
func (rm *ResolverManager) CreateResolveObject() *graphql.Object {
	fields := graphql.Fields{}

	has := false
	for name, r := range rm.resolverMap {
		has = true
		fields[name] = r.CreateField()
	}

	if !has {
		return nil
	}

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   rm.name,
			Fields: fields,
		})
}

// RegisterMutationResolver 注册函数
func (rm *ResolverManager) RegisterMutationResolver(pkg string, funcName string, funcType reflect.Type, function reflect.Value, structInstance interface{}) *RegisterInfo {
	structName := ""
	if structInstance != nil {
		structType := reflect.TypeOf(structInstance)
		structName = structType.Name()
		// structName = structType.Elem().Name()
	}

	_, ok := rm.resolverMap[funcName]
	if ok {
		panic(fmt.Errorf("Mutation Resolve [%s] exists", funcName))
	}

	info := &RegisterInfo{
		Type:    rm.name,
		Package: pkg,
		Struct:  structName,
		Func:    funcName,
	}

	r, err := rm.TryParseResolver(funcType, function, structInstance)
	if err == nil {
		rm.resolverMap[funcName] = r
		// g.mutations = append(g.mutations, info)
		// return 1
	} else {
		info.Error = err.Error()
	}

	// info.Error = err.Error()
	// g.mutations = append(g.mutations, info)

	return info
}

// func (g *GQL) registerQueryResolver(pkg string, funcName string, funcType reflect.Type, function reflect.Value, structInstance interface{}) int {
// 	// pkg := funcType.PkgPath()
// 	structName := ""
// 	if structInstance != nil {
// 		structType := reflect.TypeOf(structInstance)
// 		structName = structType.Name()
// 	}

// 	_, ok := g.queryResolvers[funcName]
// 	if ok {
// 		panic(fmt.Errorf("Query Resolve [%s] exists", funcName))
// 	}

// 	info := registerInfo{
// 		Package: pkg,
// 		Struct:  structName,
// 		Func:    funcName,
// 	}

// 	r, err := parseResolver(funcType, function, structInstance, g)
// 	if err == nil {
// 		g.queryResolvers[funcName] = r
// 		g.queries = append(g.queries, info)
// 		return 1
// 	}

// 	info.Error = err.Error()
// 	g.queries = append(g.queries, info)

// 	return 0
// }

// TryParseResolver 尝试把函数解析为 Resolver
func (rm *ResolverManager) TryParseResolver(functionType reflect.Type, function reflect.Value, structInstance interface{}) (*Resolver, error) {
	// method := reflect.TypeOf(function)
	res := &Resolver{
		manager:        rm,
		structInstance: structInstance,
		executor:       function,
	}

	// 函数返回值必须是两个，建议为两个：(其他类型, error)
	// 否则不认为是一个 resolver
	// 解析返回值
	outCount := functionType.NumOut()
	if outCount != 2 { // 返回值必须是 2 个
		return nil, errors.New("函数必须有两个返回值，且第二个必须是 error 类型")
	}
	for n := 0; n < outCount; n++ {
		outParam := functionType.Out(n)
		prop := utils.ParseTypeProp(outParam)
		if n == 0 {
			// 第一个参数
			res.out = &Field{
				Prop: prop,
			}
			// fmt.Println(res.out)
		} else if n == 1 {
			// 第二个参数 error
			res.outError = &Field{
				Prop: prop,
			}

			if prop.TypeName != "error" {
				// 第二个参数不是 error ，则确定其不是接口
				return nil, errors.New("函数的第二个返回值必须是 error 类型")
			}
		}
	}

	// 解析输入参数
	// 输入参数可以使 0~3 个，函数属于结构体，则额外多出一个参数
	// 至多包含一个 graphql.ResolveParams 参数
	// 至多包含一个 InputValidator 类型的指针参数
	// 至多包含一个 自定义 struct 类型的指针参数
	//			  struct 中的字段必须是 string, int, float, time.Time, bool 类型
	inCount := functionType.NumIn()
	for n := 0; n < inCount; n++ {
		inParam := functionType.In(n)

		prop := utils.ParseTypeProp(inParam)
		p := &Field{
			Prop: prop,
		}

		if n == 0 && structInstance != nil {
			// 结构体参数
		} else if prop.RealType == typeOfResolveParams {
			// GraphSQL param 参数,可用
			res.isGraphQLParamInParams = true
		} else if prop.RealType == typeOfInputValidator {
			// 参数验证结构
			if !prop.IsPtr {
				// 必须是指针类型
				return nil, errors.New("函数接收验证结构时，必须使用指针类型")
			}
			res.isValidatorInParams = true
		} else if prop.Kind == reflect.Struct {
			if !prop.IsPtr {
				// 必须是指针类型
				return nil, errors.New("函数接收输入的 struct 参数必须是指针类型")
			}

			// 查找注入类型，注入类型优先于参数输入类型
			inject := rm.inject.FindInject(prop.RealType)
			if inject == nil {
				// 不是注入类型，自动参数
				if res.input != nil {
					// 只能有一个 input 结构参数
					return nil, errors.New("函数只能有一个接受参数的 struct 类型")
				}

				// 结构参数，也许可用
				res.input = rm.reqObjManager.FindOrRegisterObject(p, prop.TypeName)
				if res.input == nil {
					// 不支持的结构
					return nil, errors.New("函数接收数据的 struct 的字段只能是 string, int, float, time, bool 类型")
				}
			}
		} else {
			// 不支持的数据类型
			return nil, errors.New("函数只接受 gql.InputValidator、graphql.ResolveParams 和自定义结构类型的参数")
		}

		res.params = append(res.params, p)
	}

	return res, nil
}

// CreateField 创建操作对象
func (r *Resolver) CreateField() *graphql.Field {
	var field = graphql.Field{}

	// 解析返回参数，注册 object
	if r.out.Prop.IsPrimitive {
		var err error
		field.Type, _, err = utils.TypeToGraphQLType(r.out.Prop.RealType)
		if err != nil {
			panic(err)
		}
	} else {
		// 不是原生类型
		qobj := r.manager.resObjManager.FindOrRegisterObject(r.out)
		field.Type = qobj.Object
	}

	if r.input != nil {
		iType := r.manager.reqObjManager.FindOrRegisterObject(r.input.Param, "")
		field.Args = graphql.FieldConfigArgument{
			r.input.Name: &graphql.ArgumentConfig{
				Type: iType.Object,
			},
		}
	}

	// 生成 Resolve
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		args := make([]reflect.Value, len(r.params))

		var validator *InputValidator
		var input reflect.Value
		// 构建参数
		if r.input != nil {
			// 生成输入数据结构
			var err error
			input, validator, err = parseInput(&p, r.input)
			if err != nil {
				return nil, err
			}
		}

		for n, param := range r.params {
			if r.structInstance != nil && n == 0 {
				// 结构方法，第一个参数是结构实例
				args[n] = reflect.ValueOf(r.structInstance)
				continue
			}
			if param.Prop.RealType == typeOfResolveParams {
				// 原始参数
				if param.Prop.IsPtr {
					args[n] = reflect.ValueOf(&p)
				} else {
					args[n] = reflect.ValueOf(p)
				}
			} else if param.Prop.RealType == typeOfInputValidator {
				// 验证参数
				args[n] = reflect.ValueOf(validator)
			} else if param.Prop.Kind == reflect.Struct {
				// 结构参数
				// 查找注入
				inject := r.manager.inject.FindInject(param.Prop.RealType)
				if inject != nil {
					// 注入参数
					args[n] = inject.CallFn(&p)
				} else {
					// 提交参数
					args[n] = input
				}
			}

		}

		res := r.executor.Call(args)

		resCount := len(res)
		if res == nil || resCount != 2 {
			// 没有返回值，或这返回值的数量不是两个
			return nil, nil
		}
		// 转换返回的数据
		out := res[0].Interface()
		// 转换错误信息
		errResult := res[1].Interface()
		var err error
		if errResult != nil {
			err = errResult.(error)
		}

		return out, err
	}

	return &field
}
