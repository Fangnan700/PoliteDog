package PoliteDog

import "net/http"

/*
Router 路由
*/
type Router struct {
	RouterTrie   Trie
	HandlerMap   map[string]HandlerFuc
	PreHandlers  []HandlerFuc
	PostHandlers []HandlerFuc
}

func NewRouter() *Router {
	return &Router{
		RouterTrie: Trie{
			next: &TrieNode{
				part:     "/",
				children: make([]*TrieNode, 0),
			},
		},
		HandlerMap: make(map[string]HandlerFuc),
	}
}

// 插入路由和对应handler
func (r *Router) handle(method string, pattern string, handler HandlerFuc) {
	key := md5Encode([]byte(pattern))

	r.HandlerMap[key] = handler
	r.RouterTrie.next.Insert(method, pattern, key)
}

// PreHandle 插入前置中间件
func (r *Router) PreHandle(handler ...HandlerFuc) {
	r.PreHandlers = append(r.PreHandlers, handler...)
}

// PostHandle 插入后置中间件
func (r *Router) PostHandle(handler ...HandlerFuc) {
	r.PostHandlers = append(r.PostHandlers, handler...)
}

// Use 默认使用前置中间件
func (r *Router) Use(handler ...HandlerFuc) {
	r.PreHandle(handler...)
}

func (r *Router) GET(pattern string, handler HandlerFuc) {
	r.handle(http.MethodGet, pattern, handler)
}

func (r *Router) POST(pattern string, handler HandlerFuc) {
	r.handle(http.MethodPost, pattern, handler)
}

func (r *Router) PUT(pattern string, handler HandlerFuc) {
	r.handle(http.MethodPut, pattern, handler)
}

func (r *Router) DELETE(pattern string, handler HandlerFuc) {
	r.handle(http.MethodDelete, pattern, handler)
}

/*
RouterGroup 路由组
*/
type RouterGroup struct {
	Name    string
	Routers *Router
}

func NewRouterGroup(name string) *RouterGroup {
	return &RouterGroup{
		Name:    name,
		Routers: NewRouter(),
	}
}

// 插入路由和对应handler
func (rg *RouterGroup) handle(method string, pattern string, handler HandlerFuc) {
	key := md5Encode([]byte(pattern))
	fullPattern := joinStrings(3, "/", rg.Name, pattern)

	rg.Routers.HandlerMap[key] = handler
	rg.Routers.RouterTrie.next.Insert(method, fullPattern, key)
}

// PreHandle 插入前置中间件
func (rg *RouterGroup) PreHandle(handler ...HandlerFuc) {
	rg.Routers.PreHandlers = append(rg.Routers.PreHandlers, handler...)
}

// PostHandle 插入后置中间件
func (rg *RouterGroup) PostHandle(handler ...HandlerFuc) {
	rg.Routers.PostHandlers = append(rg.Routers.PostHandlers, handler...)
}

// Use 默认使用前置中间件
func (rg *RouterGroup) Use(handler ...HandlerFuc) {
	rg.PreHandle(handler...)
}

func (rg *RouterGroup) GET(pattern string, handler HandlerFuc) {
	rg.handle(http.MethodGet, pattern, handler)
}

func (rg *RouterGroup) POST(pattern string, handler HandlerFuc) {
	rg.handle(http.MethodPost, pattern, handler)
}

func (rg *RouterGroup) PUT(pattern string, handler HandlerFuc) {
	rg.handle(http.MethodPut, pattern, handler)
}

func (rg *RouterGroup) DELETE(pattern string, handler HandlerFuc) {
	rg.handle(http.MethodDelete, pattern, handler)
}
