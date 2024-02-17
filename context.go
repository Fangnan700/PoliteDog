package PoliteDog

import (
	"net/http"
)

// Context 上下文封装
type Context struct {
	// 原始数据
	w http.ResponseWriter
	r *http.Request

	// 请求数据
	Method string
	Path   string

	// 响应数据
	Code int

	// 主体函数和中间件列表
	index    int
	handlers []HandlerFuc
}

func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Status(code int) {
	c.Code = code
	c.w.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.w.Header().Set(key, value)
}

func (c *Context) Data(data []byte) {
	c.w.Write(data)
}
