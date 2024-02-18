package PoliteDog

import (
	"github.com/fangnan700/PoliteDog/render"
	"html/template"
	"net/http"
	"net/url"
)

// Context 上下文封装
type Context struct {
	// 原始数据
	e *Dog
	w http.ResponseWriter
	r *http.Request

	// 主体函数和中间件列表
	index    int
	handlers []HandlerFuc

	// 请求数据
	Method string
	Path   string

	// 响应数据
	Code int
}

// Next 将上下文移交给下一个handler
func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

// Render 渲染器
func (c *Context) Render(w http.ResponseWriter, r render.Render) error {
	return r.Render(w)
}

// Status 返回状态码
func (c *Context) Status(code int) {
	c.Code = code
	c.w.WriteHeader(code)
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.w.Header().Set(key, value)
}

// Data 响应数据
func (c *Context) Data(code int, data []byte) error {
	c.Status(code)
	_, err := c.w.Write(data)
	return err
}

// HTML 响应HTML文本
func (c *Context) HTML(code int, html string) error {
	c.Status(code)
	return c.Render(c.w, &render.HTMLRender{
		Data:   html,
		IsTmpl: false,
	})
}

// HTMLTemplate 响应HTML模板
func (c *Context) HTMLTemplate(code int, name string, data any) error {
	c.Status(code)
	return c.Render(c.w, &render.HTMLRender{
		Name:     name,
		Data:     data,
		Template: c.e.HTMLRender.Template,
		IsTmpl:   true,
	})
}

func (c *Context) HTMLTemplateGlob(code int, name string, data interface{}, pattern string) error {
	c.Status(code)

	tmpl := template.New(name)
	tmpl, err := tmpl.ParseGlob(pattern)
	if err != nil {
		return err
	}

	err = tmpl.Execute(c.w, data)

	return err
}

// String 响应纯文本
func (c *Context) String(code int, format string, data ...any) error {
	c.Status(code)

	return c.Render(c.w, &render.StringRender{
		Format: format,
		Data:   data,
	})
}

// JSON 响应JSON数据
func (c *Context) JSON(code int, data any) error {
	c.Status(code)

	return c.Render(c.w, &render.JSONRender{
		Data: data,
	})
}

// XML 响应xml数据
func (c *Context) XML(code int, data any) error {
	c.Status(code)

	return c.Render(c.w, &render.XMLRender{
		Data: data,
	})
}

// Redirect 重定向
func (c *Context) Redirect(code int, url string) error {
	return c.Render(c.w, &render.RedirectRender{
		Code:     code,
		Request:  c.r,
		Location: url,
	})
}

// File 文件下载
func (c *Context) File(code int, filepath string) {
	c.Status(code)
	http.ServeFile(c.w, c.r, filepath)
}

// FileAttachment 自定义文件名下载
func (c *Context) FileAttachment(code int, filename string, filepath string) {
	c.Status(code)
	if isASCII(filename) {
		c.SetHeader("Content-Disposition", `attachment;filename="`+filename+`"`)
	} else {
		c.SetHeader("Content-Disposition", `attachment;filename*=UTF-8''`+url.QueryEscape(filename))
	}

	http.ServeFile(c.w, c.r, filepath)
}

// FileFromFS 从文件系统获取下载（filepath是相对于文件系统的路径）
func (c *Context) FileFromFS(code int, filepath string, fs http.FileSystem) {
	defer func(oldPath string) {
		c.r.URL.Path = oldPath
	}(c.r.URL.Path)

	c.Status(code)
	c.r.URL.Path = filepath
	http.FileServer(fs).ServeHTTP(c.w, c.r)
}
