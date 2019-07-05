package gql

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// GQL GraphSQL 辅助类
type GQL struct {
	// graphql 相关内容
	schema        *graphql.Schema
	customConfig  *handler.Config // 用户传入的 cfg
	handlerConfig *handler.Config
	// 生成的相关内容
	queries           []registerInfo
	mutations         []registerInfo
	queryResolvers    map[string]*resolver
	queryObjectMap    map[string]*queryObject
	mutationResolvers map[string]*resolver
	inputObjectMap    map[string]*inputObject
	// 注入对象结构
	inject *Inject
	// 准备要注册的函数
	tobeInject    []interface{} // 注入函数
	tobeQueries   []interface{} // 查询函数
	tobeMutations []interface{} // 操作函数
}

// 返回对象说明定义
type queryObject struct {
	object graphql.Output
	param  *param // 对应的参数
}

// 提交的参数对象定义
type inputObject struct {
	object *graphql.InputObject
	param  *param   // 对应的参数
	name   string   // 结构名称
	fields []*param // 字段列表
}

// 注册状态
type registerInfo struct {
	Package string
	Struct  string
	Func    string
	Error   string
}

var instance *GQL

func init() {
	instance = NewGQL()
}

// NewGQL 创建 GQL 实例
func NewGQL() *GQL {
	return &GQL{
		handlerConfig:     new(handler.Config),
		queryResolvers:    make(map[string]*resolver),
		queryObjectMap:    make(map[string]*queryObject),
		inputObjectMap:    make(map[string]*inputObject),
		mutationResolvers: make(map[string]*resolver),
		inject:            NewInject(),
	}
}

// Get 获取 GQL 单例实例
func Get() *GQL {
	return instance
}

// NewHandler 设置 json 是否格式化
func (g *GQL) NewHandler(cfg *handler.Config) http.Handler {
	s, err := g.GetSchema()
	if err != nil {
		panic(err)
	}

	g.customConfig = cfg
	g.handlerConfig = &handler.Config{
		Schema:           s,
		Pretty:           cfg.Pretty,
		Playground:       cfg.Playground,
		GraphiQL:         cfg.GraphiQL,
		ResultCallbackFn: cfg.ResultCallbackFn,
		FormatErrorFn:    cfg.FormatErrorFn,
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			var root map[string]interface{}
			if cfg.RootObjectFn != nil {
				root = cfg.RootObjectFn(ctx, r)
			}
			if root == nil {
				root = make(map[string]interface{})
			}
			// 存储 context 和 request
			g.inject.StoreContext(ctx, r, root)
			return root
		},
	}

	// cfg.Schema = &s

	return handler.New(g.handlerConfig)
}

// GetSchema 获取 GraphQL 结构，如果哈没有创建，则创建
func (g *GQL) GetSchema() (*graphql.Schema, error) {
	if g.schema != nil {
		return g.schema, nil
	}

	// 注册注入函数
	for _, fn := range g.tobeInject {
		g.inject.Inject(fn)
	}
	// 注册查询函数
	for _, obj := range g.tobeQueries {
		g.doRegisterQuery(obj)
	}
	// 注册操作函数
	for _, obj := range g.tobeMutations {
		g.doRegisterMutation(obj)
	}

	// 生成 graphql 结构
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.GetQuery(),
			Mutation: g.GetMutation(),
		},
	)
	g.schema = &schema

	return g.schema, err
}

func (g *GQL) createSummary(filter func(info *registerInfo) bool) string {
	query := "Query:"
	for _, info := range g.queries {
		if filter != nil && !filter(&info) {
			continue
		}
		query += "\n\t"
		if info.Struct != "" {
			query += info.Struct + "." + info.Func + "@" + info.Package
		} else {
			query += info.Func + "@" + info.Package
		}
		if info.Error == "" {
			query += "\t[OK]"
		} else {
			query += "\t[" + info.Error + "]"
		}
	}

	mutation := "\nMutation:"
	for _, info := range g.mutations {
		if filter != nil && !filter(&info) {
			continue
		}
		mutation += "\n\t"
		if info.Struct != "" {
			mutation += info.Struct + "." + info.Func + "@" + info.Package
		} else {
			mutation += info.Func + "@" + info.Package
		}
		if info.Error == "" {
			mutation += "\t[OK]"
		} else {
			mutation += "\t[" + info.Error + "]"
		}
	}

	return query + mutation + "\n"
}

// Summary 注册说明
func (g *GQL) Summary() string {
	return "** Register List **\n---------------------------------------------------------------------------------\n" +
		g.createSummary(nil) +
		"---------------------------------------------------------------------------------\n"
}

// SummaryOfFailed 列举注册失败的函数信息
func (g *GQL) SummaryOfFailed() string {
	return "** Register Failed List **\n---------------------------------------------------------------------------------\n" +
		g.createSummary(func(info *registerInfo) bool {
			return info.Error != ""
		}) +
		"---------------------------------------------------------------------------------\n"
}

// GetQuery 获取查询类型
func (g *GQL) GetQuery() *graphql.Object {
	fields := graphql.Fields{}

	has := false
	for name, r := range g.queryResolvers {
		has = true
		fields[name] = r.GetQueryField()
	}

	if !has {
		return nil
	}

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: fields,
		})
}

// GetMutation 获取操作类型
func (g *GQL) GetMutation() *graphql.Object {
	fields := graphql.Fields{}

	has := false
	for name, r := range g.mutationResolvers {
		has = true
		fields[name] = r.GetQueryField()
	}
	if !has {
		return nil
	}
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: fields,
		})
}

// RegisterMutation 注册操作，加入到待注册列表
func (g *GQL) RegisterMutation(query interface{}) {
	g.tobeMutations = append(g.tobeMutations, query)
}

// 执行注册操作
func (g *GQL) doRegisterMutation(mutation interface{}) int {
	mutationType := reflect.TypeOf(mutation)
	kind := mutationType.Kind()
	// pkg := mutationType.PkgPath()
	if kind == reflect.Func {
		// 是一个函数
		val := reflect.ValueOf(mutation)
		funcInfo := ParseFuncInfo(mutation)
		// name := getFuncName(mutation, '/', '.')
		// pkg := getFuncPackage(mutation)
		return g.registerMutationResolver(funcInfo.Pkg, funcInfo.Name, mutationType, val, nil)
	} else if kind == reflect.Struct {
		// 是一个结构，遍历其所有函数
		sum := 0
		pkg := mutationType.PkgPath()
		for n := 0; n < mutationType.NumMethod(); n++ {
			method := mutationType.Method(n)
			name := method.Name
			sum += g.registerMutationResolver(pkg, name, method.Type, method.Func, mutation)
		}
		return sum
	}

	return 0
}

// RegisterInject 注册注入函数,只是加入到待注册列表
func (g *GQL) RegisterInject(injectFn interface{}) {
	g.tobeInject = append(g.tobeInject, injectFn)
}

// var typeOfContext = reflect.TypeOf(context.TODO())

// RegisterQuery 注册查询
func (g *GQL) RegisterQuery(query interface{}) {
	g.tobeQueries = append(g.tobeQueries, query)
}

// 执行注册查询
// @return 返回注册成功的函数数量
func (g *GQL) doRegisterQuery(query interface{}) int {
	queryType := reflect.TypeOf(query)
	kind := queryType.Kind()
	if kind == reflect.Func {
		// 是一个函数
		val := reflect.ValueOf(query)
		funcInfo := ParseFuncInfo(query)
		// name := getFuncName(query, '/', '.')
		// pkg := getFuncPackage(query)
		return g.registerQueryResolver(funcInfo.Pkg, funcInfo.Name, queryType, val, nil)
	} else if kind == reflect.Struct {
		// 是一个结构，遍历其所有函数
		sum := 0
		pkg := queryType.PkgPath()
		for n := 0; n < queryType.NumMethod(); n++ {
			method := queryType.Method(n)
			name := method.Name
			sum += g.registerQueryResolver(pkg, name, method.Type, method.Func, query)
		}
		return sum
	}
	return 0
}

func (g *GQL) registerQueryResolver(pkg string, funcName string, funcType reflect.Type, function reflect.Value, structInstance interface{}) int {
	// pkg := funcType.PkgPath()
	structName := ""
	if structInstance != nil {
		structType := reflect.TypeOf(structInstance)
		structName = structType.Name()
	}

	_, ok := g.queryResolvers[funcName]
	if ok {
		panic(fmt.Errorf("Query Resolve [%s] exists", funcName))
	}

	info := registerInfo{
		Package: pkg,
		Struct:  structName,
		Func:    funcName,
	}

	r, err := parseResolver(funcType, function, structInstance, g)
	if err == nil {
		g.queryResolvers[funcName] = r
		g.queries = append(g.queries, info)
		return 1
	}

	info.Error = err.Error()
	g.queries = append(g.queries, info)

	return 0
}

func (g *GQL) registerMutationResolver(pkg string, funcName string, funcType reflect.Type, function reflect.Value, structInstance interface{}) int {
	// pkg := funcType.PkgPath()
	structName := ""
	if structInstance != nil {
		structType := reflect.TypeOf(structInstance)
		structName = structType.Elem().Name()
	}

	_, ok := g.mutationResolvers[funcName]
	if ok {
		panic(fmt.Errorf("Mutation Resolve [%s] exists", funcName))
	}

	info := registerInfo{
		Package: pkg,
		Struct:  structName,
		Func:    funcName,
	}

	r, err := parseResolver(funcType, function, structInstance, g)
	if err == nil {
		g.mutationResolvers[funcName] = r
		g.mutations = append(g.mutations, info)
		return 1
	}

	info.Error = err.Error()
	g.mutations = append(g.mutations, info)

	return 0
}

// 查找查询对象，如果找不到则注册
func (g *GQL) findOrRegisterQueryObject(field *param) *queryObject {
	list := field.Prop.IsList
	key, keyList := field.Prop.Key(), field.Prop.KeyList()
	var obj *queryObject
	var ok bool
	if list {
		// 列表
		obj, ok = g.queryObjectMap[keyList]
		if !ok {
			// 没有列表
			obj, ok = g.queryObjectMap[key]
			if ok {
				// 有单独的注册对象
				// 注册列表即可
				obj = g.registerList(field, graphql.NewList(obj.object))
				return obj
			}
		}
	} else {
		obj, ok = g.queryObjectMap[key]
	}
	if !ok {
		// 没有找到，注册
		objFields := graphql.Fields{}

		p := field.Prop
		name := p.TypeName
		// typeField := new(graphql.Field)
		for n := 0; n < p.Type.NumField(); n++ {
			field := p.Type.Field(n)
			typeField := new(graphql.Field)
			ftype, isStruct := parseGtraphType(&field)
			if isStruct {
				structFieldProp := parseProp(field.Type)
				structFieldParam := &param{Prop: structFieldProp}
				fobj := g.findOrRegisterQueryObject(structFieldParam)
				typeField.Type = fobj.object
			} else {
				typeField.Type = ftype
			}
			id := parseFieldName(&field)
			objFields[id] = typeField
		}

		// 注册单个查询对象
		gobj := graphql.NewObject(graphql.ObjectConfig{
			Name:   name,
			Fields: objFields,
		})

		obj = g.registerObject(field, gobj)
		if list {
			// 是数组
			obj = g.registerList(field, graphql.NewList(gobj))
		}
	}
	return obj
}

// 查找查询对象，如果找不到则注册
func (g *GQL) findOrRegisterInputObject(field *param, name string) *inputObject {
	list := field.Prop.IsList
	if list {
		// 如果是列表，则返回空，不可以使用列表提交
		return nil
	}
	key := field.Prop.Key()

	obj, ok := g.inputObjectMap[key]

	if !ok {
		if name == "" {
			// name 为空时，必须找到对象
			panic(fmt.Errorf("Cann't find input object %s", field))
		}
		// 没有找到，注册
		objFields := graphql.InputObjectConfigFieldMap{}

		p := field.Prop
		// name := p.TypeName
		var fields []*param
		for n := 0; n < p.Type.NumField(); n++ {
			field := p.Type.Field(n)
			typeField := new(graphql.InputObjectFieldConfig)
			ftype, isStruct := parseGtraphType(&field)
			if isStruct {
				// 是结构类型，不支持
				return nil
			}
			id := parseFieldName(&field)
			prop := parseProp(field.Type)
			// 存储结构字段
			prop.structField = field
			fields = append(fields, &param{
				Name:     field.Name,
				JSONName: id,
				Prop:     prop,
			})

			typeField.Type = ftype

			objFields[id] = typeField
		}

		// 注册单个查询对象
		gobj := graphql.NewInputObject(
			graphql.InputObjectConfig{
				Name:   "input" + name,
				Fields: objFields,
			})

		obj = &inputObject{
			object: gobj,
			param:  field,
			name:   name,
			fields: fields,
		}
		g.inputObjectMap[key] = obj
	}
	return obj
}

func (g *GQL) registerObject(param *param, obj *graphql.Object) *queryObject {
	qobj := &queryObject{
		object: obj,
		param:  param,
	}
	g.queryObjectMap[param.Prop.Key()] = qobj
	return qobj
}

func (g *GQL) registerList(param *param, list *graphql.List) *queryObject {
	qobj := &queryObject{
		object: list,
		param:  param,
	}
	g.queryObjectMap[param.Prop.KeyList()] = qobj
	return qobj
}

func parseGtraphType(field *reflect.StructField) (grapghType graphql.Output, isStruct bool) {
	tp := field.Type

	kind := tp.Kind()

	isPtr := kind == reflect.Ptr
	isList := kind == reflect.Slice

	if isPtr || isList {
		// 指针、切片、数组
		// field.Type.Elem()
		kind = field.Type.Elem().Kind()
		tp = field.Type.Elem()
	}

	gType, isStruct, err := TypeToGraphQLType(tp)
	if err != nil {
		panic(err)
	}
	if isStruct {
		return nil, true
	}
	return gType, false
	// if IsIntType(kind) {
	// 	return graphql.Int, false
	// }
	// if IsFloatType(kind) {
	// 	return graphql.Float, false
	// }
	// if IsStringType(kind) {
	// 	return graphql.String, false
	// }
	// if IsBoolType(kind) {
	// 	return graphql.Boolean, false
	// }
	// if IsTimeType(tp) {
	// 	return graphql.DateTime, false
	// }
	// if IsStructType(kind) {
	// 	return nil, true
	// }

	// panic(fmt.Errorf("不支持 %s 类型", tp.Name()))
}

func parseFieldName(field *reflect.StructField) string {
	name := field.Tag.Get("json")
	if name == "" {
		name = field.Name
		tmp := strings.ToLower(string(name[0]))
		name = strings.Replace(name, string(name[0]), tmp, 1)
	}
	return name
}
