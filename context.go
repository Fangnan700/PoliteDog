package PoliteDog

import (
	"errors"
	"github.com/fangnan700/PoliteDog/binding"
	"github.com/fangnan700/PoliteDog/render"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

const defaultMultipartMaxMemory = 32

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
	Method     string
	Path       string
	queryCache url.Values
	formCache  url.Values

	// 响应数据
	Code int

	// 其它参数
	DisallowUnknownFields bool // 是否校验json对应结构体字段
}

/**
参数解析
*/

// 初始化queryCache
func (c *Context) initQueryCache() {
	if c.r != nil {
		c.queryCache = c.r.URL.Query()
	} else {
		c.queryCache = url.Values{}
	}
}

// GetQuery 获取query参数
func (c *Context) GetQuery(key string) any {
	c.initQueryCache()
	return c.queryCache.Get(key)
}

// GetQueryArray 获取query参数切片
func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQueryCache()
	val, ok := c.queryCache[key]
	return val, ok
}

// 初始化PostFormCache
func (c *Context) initPostFormCache() {
	if c.r != nil {
		err := c.r.ParseMultipartForm(defaultMultipartMaxMemory)
		if err != nil {
			// 这里由于接收的是通用表单，所以忽略ErrNotMultipart
			if !errors.Is(err, http.ErrNotMultipart) {
				c.e.logger.Error(err.Error())
			}
		}
		c.formCache = c.r.PostForm
	} else {
		c.formCache = url.Values{}
	}
}

// GetMultipartForm 获取原始MultipartForm
func (c *Context) GetMultipartForm() (*multipart.Form, error) {
	err := c.r.ParseMultipartForm(defaultMultipartMaxMemory)
	return c.r.MultipartForm, err
}

// GetPostForm 获取postForm
func (c *Context) GetPostForm(key string) any {
	c.initPostFormCache()
	return c.formCache.Get(key)
}

// GetPostFormArray 获取postForm切片
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.initPostFormCache()
	val, ok := c.formCache[key]
	return val, ok
}

// GetFormFile 获取表单文件，返回文件头
func (c *Context) GetFormFile(key string) (*multipart.FileHeader, error) {
	_, header, err := c.r.FormFile(key)
	return header, err
}

// SaveUploadFile 封装文件上传并保存的方法
func (c *Context) SaveUploadFile(fileHeader *multipart.FileHeader, savePath string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

/**
数据绑定
*/

func (c *Context) MustBindWith(obj any, binding binding.Binding) error {
	return binding.Bind(c.r, obj)
}

// BindJSON 解析JSON参数
func (c *Context) BindJSON(obj any) error {
	jsonBind := binding.JSONBind
	jsonBind.DisallowUnknownFields = c.DisallowUnknownFields
	return c.MustBindWith(obj, jsonBind)
}

// BindXML 解析XML参数
func (c *Context) BindXML(obj any) error {
	xmlBind := binding.XMLBind
	return c.MustBindWith(obj, xmlBind)
}

/**
上下文操作
*/

// Next 将上下文移交给下一个handler
func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

// Abort 将上下文移交给最后一个handler
func (c *Context) Abort() {
	c.index = len(c.handlers) - 1
	c.handlers[c.index](c)
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

/**
数据响应、渲染
*/

// Render 渲染器
func (c *Context) Render(w http.ResponseWriter, r render.Render) error {
	return r.Render(w)
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
