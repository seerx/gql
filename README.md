# gql
简化 graphql-go 开发，使用 graphql 的原始接口开发非常繁琐，稍有不慎就会出错。该项目就是为了简化 golang 开发 graphql 过程而生。

# 简化内容
<ol>
<li>省去定义 graphql 对象的过程</li>
<li>省去定义 query 和 mutation 的过程</li>
<li>省去手动获取请求参数的过程，并提供参数验证功能</li>
<li>提供注入功能，简化登录验证等功能</li>
</ol>

# 例子
<pre>
// Goods 商品信息
type Goods struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	URL   string    `json:"url"`
	Time  time.Time `json:"time"`
}

// GoodsList 查询函数定义
func GoodsList() ([]Goods, error) {
	return []entities.Goods{
		entities.Goods{
			ID: "A1", Name: "A-test1",
		},
		entities.Goods{
			ID: "A2", Name: "A-test2",
		},
	}, nil
}

func init() {
	// 注册查询
	gql.Get().RegisterQuery(GoodsList)
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
