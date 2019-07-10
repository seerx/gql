package gql

import (
	"context"
	"net/http"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/seerx/gql/pkg/gqlh"
	"github.com/seerx/gql/utils"
)

// GQL GraphSQL 辅助类
type GQL struct {
	// graphql 相关内容
	schema        *graphql.Schema
	customConfig  *handler.Config // 用户传入的 cfg
	handlerConfig *handler.Config
	// 生成的相关内容
	registerInfos []*gqlh.RegisterInfo

	queryManager          *gqlh.ResolverManager
	mutationManager       *gqlh.ResolverManager
	responseObjectManager *gqlh.ResponseObjectManager
	requestObjectManager  *gqlh.RequestObjectManager
	// 注入对象结构
	inject *gqlh.Inject
	// 准备要注册的函数
	tobeInject    []interface{} // 注入函数
	tobeQueries   []interface{} // 查询函数
	tobeMutations []interface{} // 操作函数
}

var instance *GQL

func init() {
	instance = NewGQL()
}

// NewGQL 创建 GQL 实例
func NewGQL() *GQL {
	inject := gqlh.NewInject()
	resObj := gqlh.NewResponseObjectManager()
	reqObj := gqlh.NewRequestObjectManager()
	return &GQL{
		handlerConfig:         new(handler.Config),
		responseObjectManager: resObj, //  gqlh.NewResponseObjectManager(),
		requestObjectManager:  reqObj, // gqlh.NewRequestObjectManager(),
		queryManager:          gqlh.NewQueryResolverManager(inject, resObj, reqObj),
		mutationManager:       gqlh.NewMutationResolverManager(inject, resObj, reqObj),
		inject:                inject, //   gqlh.NewInject(),
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
		// g.doRegisterQuery(obj)
		g.doRegisterResolver(g.queryManager, obj)
	}
	// 注册操作函数
	for _, obj := range g.tobeMutations {
		// g.doRegisterMutation(obj)
		g.doRegisterResolver(g.mutationManager, obj)
	}

	// 生成 graphql 结构
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.queryManager.CreateResolveObject(),
			Mutation: g.mutationManager.CreateResolveObject(),
		},
	)
	if err == nil {
		g.schema = &schema
	}
	return g.schema, err
}

func (g *GQL) createSummary(filter func(info *gqlh.RegisterInfo) bool) string {
	query := "Query:"
	mutation := "\nMutation:"
	for _, info := range g.registerInfos {
		if filter != nil && !filter(info) {
			continue
		}
		str := "\n\t"
		if info.Struct != "" {
			str += info.Struct + "." + info.Func + "@" + info.Package
		} else {
			str += info.Func + "@" + info.Package
		}
		if info.Error == "" {
			str += "\t[OK]"
		} else {
			str += "\t[" + info.Error + "]"
		}
		if info.Type == "Query" {
			query += str
		} else {
			mutation += str
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
		g.createSummary(func(info *gqlh.RegisterInfo) bool {
			return info.Error != ""
		}) +
		"---------------------------------------------------------------------------------\n"
}

// RegisterMutation 注册操作，加入到待注册列表
// @param mutation 可以是一个函数，也可以是一个有多个函数的结构体
func (g *GQL) RegisterMutation(mutation interface{}) {
	g.tobeMutations = append(g.tobeMutations, mutation)
}

func (g *GQL) doRegisterResolver(manager *gqlh.ResolverManager, resolveFunc interface{}) {
	funcType := reflect.TypeOf(resolveFunc)
	kind := funcType.Kind()
	if kind == reflect.Func {
		// 是一个函数
		val := reflect.ValueOf(resolveFunc)
		funcInfo := utils.ParseFuncInfo(resolveFunc)
		info := manager.RegisterMutationResolver(funcInfo.Pkg, funcInfo.Name, funcType, val, nil)
		g.registerInfos = append(g.registerInfos, info)
	} else if kind == reflect.Struct {
		// 是一个结构，遍历其所有函数
		// sum := 0
		pkg := funcType.PkgPath()
		for n := 0; n < funcType.NumMethod(); n++ {
			method := funcType.Method(n)
			name := method.Name
			info := manager.RegisterMutationResolver(pkg, name, method.Type, method.Func, resolveFunc)
			g.registerInfos = append(g.registerInfos, info)
		}
	}
}

// 执行注册操作
// func (g *GQL) doRegisterMutation(mutation interface{}) {
// 	g.doRegisterResolver(g.mutationManager, mutation)
// }

// RegisterInject 注册注入函数,只是加入到待注册列表
// @param injectFn 必须是一个函数，且函数必须是一下形式
// 		func injectFn(ctx, context.Context, r *http.Request) *CustomStruct
func (g *GQL) RegisterInject(injectFn interface{}) {
	g.tobeInject = append(g.tobeInject, injectFn)
}

// RegisterQuery 注册查询
// @param query 可以是一个函数，也可以是一个有多个函数的结构体
func (g *GQL) RegisterQuery(query interface{}) {
	g.tobeQueries = append(g.tobeQueries, query)
}

// 执行注册查询
// @return 返回注册成功的函数数量
// func (g *GQL) doRegisterQuery(query interface{}) {
// 	g.doRegisterResolver(g.queryManager, query)
// }
