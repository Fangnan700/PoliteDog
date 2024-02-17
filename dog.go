package PoliteDog

import (
	"fmt"
	"log"
	"net/http"
)

// Dog 核心引擎结构体
type Dog struct {
	Routers      []*Router
	RouterGroups []*RouterGroup
}

func NewDog() *Dog {
	return &Dog{
		Routers: make([]*Router, 0),
	}
}

// RegisterRouters 将路由注册到引擎
func (dog *Dog) RegisterRouters(routers ...*Router) {
	for _, r := range routers {
		dog.Routers = append(dog.Routers, r)
	}
}

// RegisterRouterGroup 解析路由组，将路由注册到引擎
func (dog *Dog) RegisterRouterGroup(groups ...*RouterGroup) {
	for _, group := range groups {
		dog.Routers = append(dog.Routers, group.Routers)
		dog.RouterGroups = append(dog.RouterGroups, group)
	}
}

// ServeHTTP 预处理Http请求
func (dog *Dog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		w:        w,
		r:        r,
		Method:   r.Method,
		index:    -1,
		handlers: make([]HandlerFuc, 0),
	}

	path := r.URL.Path
	matched := false
	for _, router := range dog.Routers {
		trieNode := router.RouterTrie.next.Search(path)

		// 匹配到路由
		if trieNode != nil && trieNode.end {
			if trieNode.method != ctx.Method {
				ctx.Status(http.StatusMethodNotAllowed)
				return
			}
			matched = true

			// 根据key提取handler
			key := trieNode.key
			handle := router.HandlerMap[key]

			// 提取中间件，和handle一起注册到上下文
			ctx.handlers = append(ctx.handlers, router.PreHandlers...)
			ctx.handlers = append(ctx.handlers, handle)
			ctx.handlers = append(ctx.handlers, router.PostHandlers...)

			// 将上下文移交给对应的handler
			ctx.Next()
		}
	}

	// 未找到路由
	if !matched {
		ctx.Status(http.StatusNotFound)
	}
}

// Run 启动！
func (dog *Dog) Run(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)

	err := http.ListenAndServe(addr, dog)
	if err != nil {
		log.Fatalln(err)
	}
}