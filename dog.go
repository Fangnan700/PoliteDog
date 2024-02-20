package PoliteDog

import (
	"fmt"
	"github.com/fangnan700/PoliteDog/logger"
	"github.com/fangnan700/PoliteDog/render"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Dog 核心引擎结构体
type Dog struct {
	pool         sync.Pool
	logger       *logger.Logger
	Routers      []*Router
	RouterGroups []*RouterGroup
	Middlewares  []HandlerFuc
	TmplFuncMap  template.FuncMap
	HTMLRender   render.HTMLRender
}

func NewDog() *Dog {
	dog := &Dog{
		Routers: make([]*Router, 0),
	}
	dog.pool.New = func() any {
		return dog.allocateContext()
	}
	dog.logger = logger.DefaultLogger()

	return dog
}

// 通过同步对象池来解决context频繁创建的问题
func (dog *Dog) allocateContext() any {
	return &Context{
		e: dog,
	}
}

// SetFuncMap 设置模板渲染过程中可能使用的自定义函数
func (dog *Dog) SetFuncMap(funcMap template.FuncMap) {
	dog.TmplFuncMap = funcMap
}

// SetTemplate 允许开发者自己设置模板
func (dog *Dog) SetTemplate(tmpl *template.Template) {
	dog.HTMLRender = render.HTMLRender{Template: tmpl}
}

// LoadTemplate 加载模板
func (dog *Dog) LoadTemplate(pattern string) {
	tmpl := template.Must(template.New("").Funcs(dog.TmplFuncMap).ParseGlob(pattern))
	dog.SetTemplate(tmpl)
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

// ServeHTTP
func (dog *Dog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := dog.pool.Get().(*Context)
	ctx.w = w
	ctx.r = r
	ctx.Method = r.Method
	ctx.Path = r.URL.Path
	ctx.index = -1
	ctx.Code = 0
	ctx.handlers = make([]HandlerFuc, 0)

	dog.HttpRequestHandler(ctx)
	dog.pool.Put(ctx)
}

// HttpRequestHandler 预处理Http请求
func (dog *Dog) HttpRequestHandler(ctx *Context) {
	path := ctx.r.URL.Path
	matched := false
	methodHit := false

	// 注册异常捕获中间件
	ctx.handlers = append(ctx.handlers, Recovery)

	for _, router := range dog.Routers {
		trieNode := router.RouterTrie.next.Search(path)

		// 匹配到路由
		if trieNode != nil && trieNode.end {
			matched = true

			// 根据key提取handler
			key := trieNode.key
			handle := router.HandlerMap[key]

			// 提取中间件，和handle一起注册到上下文
			ctx.handlers = append(ctx.handlers, router.PreHandlers...)
			ctx.handlers = append(ctx.handlers, handle)
			ctx.handlers = append(ctx.handlers, router.PostHandlers...)

			// 校验请求方法
			if trieNode.method == ctx.Method {
				methodHit = true
			}
		}
	}

	// 注册日志中间件
	ctx.handlers = append(ctx.handlers, dog.logReq)

	if matched {
		if methodHit {
			ctx.Next()
			return
		} else {
			ctx.Data(http.StatusMethodNotAllowed, nil)
			ctx.Abort()
			return
		}
	} else {
		ctx.Status(http.StatusNotFound)
		ctx.Abort()
		return
	}
}

func init() {
	clearTerminal()
	fmt.Println(
		"\u001B[36m" +
			"\t _____          _  _  _          _____\n" +
			"\t|  __ \\        | |(_)| |        |  __ \\\n" +
			"\t| |__) |  ___  | | _ | |_   ___ | |  | |  ___    __ _\n" +
			"\t|  ___/  / _ \\ | || || __| / _ \\| |  | | / _ \\  / _` |\n" +
			"\t| |     | (_) || || || |_ |  __/| |__| || (_) || (_| |\n" +
			"\t|_|      \\___/ |_||_| \\__| \\___||_____/  \\___/  \\__, |\n" +
			"\t                                                 __/ |\n" +
			"\t                                                |___/\n",
	)
}

// SetLogPath 设置日志路径
func (dog *Dog) SetLogPath(logPath string) {
	dog.logger.SetLogPath(logPath)
}

// 打印请求日志
func (dog *Dog) logReq(ctx *Context) {
	msg := fmt.Sprintf("%3d %-8s %s", ctx.Code, ctx.Method, ctx.Path)
	if ctx.Code == http.StatusOK {
		dog.logger.Info(msg)
	} else {
		dog.logger.Warning(msg)
	}
}

// Run 启动！
func (dog *Dog) Run(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	dog.logger.Info(fmt.Sprintf("PoliteDog running at: %s", addr))

	err := http.ListenAndServe(addr, dog)
	if err != nil {
		log.Fatalln(err)
	}
}
