package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/seerx/gql"
)

func init() {
	gql.Get().RegisterInject(InjectAccount)
	gql.Get().RegisterInject(InjectInterface)
}

// Account 测试
type Account struct {
	ID   string
	Name string
}

func (a *Account) String() string {
	return fmt.Sprintf("{id: %s, name: %s}", a.ID, a.Name)
}

type A interface {
	Do(a string)
}

type AI struct {
}

func (*AI) Do(a string) {
	fmt.Println(a)
}

func (*AI) Close() error {
	fmt.Println("Call Close")
	return nil
}

// InjectInterface 注入接口
func InjectInterface(ctx context.Context, r *http.Request) A {
	return &AI{}
}

// InjectAccount 测试
func InjectAccount(ctx context.Context, r *http.Request) *Account {
	// 可以读取 http.Request 的 cookie 等信息，完后从缓存（或数据库）中获取 Account 信息
	// 如果无法获取，则可以 panic 一个错误信息，此错误将会返回到客户端
	// if Cann't get account info {
	// 	panic(errors.New("Need Login"))
	// }
	return &Account{
		ID:   "account.id",
		Name: "Init by InjectAccount ",
	}
}
