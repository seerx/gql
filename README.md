# 该项目不再维护
该包的功能在 goql 中被继承，请跳转到 <a href="https://github.com/seerx/goql">goql</a>

# gql
简化 <a href="https://github.com/graphql-go/graphql">graphql-go</a> 开发，graphql-ql 提供的原始接口开发非常繁琐，稍有不慎就会出错。该项目就是为了简化 golang 开发 graphql 过程，让开发人员把更多的精力投入到业务逻辑的开发中。

# 简化内容
<ol>
<li>省去定义 graphql 对象的过程，只需定义对应的结构体即可</li>
<li>省去定义 query 和 mutation 的过程，只需定义 go 函数即可</li>
<li>省去手动获取请求参数的过程，并提供参数验证功能（完善中）</li>
<li>提供注入功能，简化登录验证等功能</li>
</ol>

# 例子
<pre>
// goods 商品信息
type goods struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	URL   string    `json:"url"`
	Time  time.Time `json:"time"`
}

// goodsList 查询函数定义
func goodsList() ([]goods, error) {
	return []goods{
		goods{
			ID: "A1", Name: "A-test1",
		},
		goods{
			ID: "A2", Name: "A-test2",
		},
	}, nil
}

func init() {
	// 注册查询
	gql.Get().RegisterQuery(goodsList)
}

func main() {
	// 在 8080 端口启动服务
	g := gql.Get()
	handler := g.NewHandler(&handler.Config{
		Pretty:   true,
		GraphiQL: true,
	})

	fmt.Print(g.Summary())

	http.Handle("/graphql", handler)
	fmt.Println("The api server will run on port : ", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil) 
}
</pre>

启动后在浏览器输入 http://localhost:8080/graphql 即可看到接口信息

# 约定
<ol>
	<li>定义的对象名称与结构名称完全一致</li>
	<li>输入参数只允许使用自定义结构，且只能有一个输入参数，输入参的参数名称以结构名称定义，graphql 对象类型将以 “input” 打头加上结构名称定义</li>
	<li>输入结构允许嵌套，但是不允许出现匿名结构</li>
	<li>输入参数如果是自定义结构，要使用指针形式</li>
	<li>输入参数不能出现切片、数组、映射等类型（同时提交多条记录，可以使用别名）</li>
	<li>只有返回值是一个自定义结构（指针）和 error 的函数，且输入参数为一个自定义结构（指针）、任意个数注入结构（指针）、*gqlh.InputValidator 的组合才能作为 Query 和 Mutation 对象</li>
	<li>定义的 Query 和 Mutation 名称与函数名称完全一致</li>
	<li>出现 Query 或 Mutation 的函数名称相同时，将舍弃后面的函数</li>
	<li>注入函数必须是固定形式</li>
</ol>

# 更详细的使用方法 examples
<ol>
	<li><a href="https://github.com/seerx/gql/tree/master/examples/hello">hello 简单示例</a></li>
	<li><a href="https://github.com/seerx/gql/tree/master/examples/query">query Query 示例</a></li>
	<li><a href="https://github.com/seerx/gql/tree/master/examples/mutation">mutation Mutation 示例</a></li>
	<li><a href="https://github.com/seerx/gql/tree/master/examples/inject">inject 注入示例</a></li>
	<li><a href="https://github.com/seerx/gql/tree/master/examples/checkparams">checkparams 参数检测示例</a></li>
</ol>
