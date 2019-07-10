# gql
简化 <a href="https://github.com/graphql-go/graphql">graphql-go</a> 开发，graphql-ql 提供的原始接口开发非常繁琐，稍有不慎就会出错。该项目就是为了简化 golang 开发 graphql 过程而生。

# 简化内容
<ol>
<li>省去定义 graphql 对象的过程</li>
<li>省去定义 query 和 mutation 的过程</li>
<li>省去手动获取请求参数的过程，并提供参数验证功能</li>
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
	<li>如果同一个结构既用做输入参数，那么定义的对象将以 “input” 打头</li>
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
