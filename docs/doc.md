# 快速使用PoliteDog



## API示例

### 路由

目前已支持：

- GET
- POST
- PUT
- DELETE



**简单示例**

```go
import "github.com/fangnan700/PoliteDog"

func main() {
	// 创建引擎
	dog := PoliteDog.NewDog()

    // 创建路由
	router := PoliteDog.NewRouter()
	router.GET("/user/info", func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("方法：GET 访问路由：/user/info"))
	})
	router.POST("/user/login", func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("方法：POST 访问路由：/user/login"))
	})

	// 将路由注册到引擎
	dog.RegisterRouters(router)

	// 启动引擎，需指定监听地址和端口
	dog.Run("127.0.0.1", 8080)
}
```

**同时注册多个路由**

```go
// 创建路由
router1 := PoliteDog.NewRouter()
router1.GET("/user/info", func(ctx *PoliteDog.Context) {
	ctx.Data(http.StatusOK, []byte("方法：GET 访问路由：/user/info"))
})

router2 := PoliteDog.NewRouter()
router2.POST("/user/login", func(ctx *PoliteDog.Context) {
	ctx.Data(http.StatusOK, []byte("方法：POST 访问路由：/user/login"))
})

// 将路由注册到引擎
dog.RegisterRouters(router1, router2)
```





**RESTful API**

```go
import "github.com/fangnan700/PoliteDog"

func main() {
	// 创建引擎
	dog := PoliteDog.NewDog()

    // 创建路由
	router := PoliteDog.NewRouter()
	router.PUT("/user/info/:id", func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("方法：PUT 访问路由：/user/info/:id"))
	})
	router.DELETE("/user/delete/*", func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("方法：DELETE 访问路由：/user/delete/*"))
	})

	// 将路由注册到引擎
	dog.RegisterRouters(router)

	// 启动引擎，需指定监听地址和端口
	dog.Run("127.0.0.1", 8080)
}
```





### 中间件

```go
import "github.com/fangnan700/PoliteDog"

func main() {
	// 创建引擎
	dog := PoliteDog.NewDog()

	// 创建路由
	router := PoliteDog.NewRouter()
	router.GET("/user/info", func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("方法：GET 访问路由：/user/info\n"))
	})

	// 注册前置中间件
	router.PreHandle(func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("前置中间件正在工作\n"))
	})

	// 注册后置中间件
	router.PostHandle(func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("后置中间件正在工作\n"))
	})

	// 将路由注册到引擎
	dog.RegisterRouters(router)

	// 启动引擎，需指定监听地址和端口
	dog.Run("127.0.0.1", 8080)
}
```

请求处理顺序为：

```text
前置中间件 -> 主体函数 -> 后置中间件
```

执行上述代码，请求后应得到响应：

![image-20240217112600242](https://yvling-typora-image-1257337367.cos.ap-nanjing.myqcloud.com/typora/image-20240217112600242.png)



当然，也可以直接使用 `router.Use()` 来注册中间件，默认注册的是前置中间件：

```go
router.Use(func(ctx *PoliteDog.Context) {
	ctx.Data(http.StatusOK, []byte("默认中间件正在工作\n"))
})
```

请求后应得到响应：

![image-20240217112753297](https://yvling-typora-image-1257337367.cos.ap-nanjing.myqcloud.com/typora/image-20240217112753297.png)



### 路由组

**路由组的大多数使用方法与基本路由一致，可直接参考基本路由。**

```go
import "github.com/fangnan700/PoliteDog"

func main() {
	// 创建引擎
	dog := PoliteDog.NewDog()

	// 创建路由组
	group := PoliteDog.NewRouterGroup("admin")
	group.GET("/info", func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("方法：GET 访问路由：" + ctx.Path + "\n"))
	})

	// 注册前置中间件
	group.PreHandle(func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("前置中间件正在工作\n"))
	})

	// 注册后置中间件
	group.PostHandle(func(ctx *PoliteDog.Context) {
		ctx.Data(http.StatusOK, []byte("后置中间件正在工作\n"))
	})

	// 将路由组注册到引擎
	dog.RegisterRouterGroup(group)

	// 启动引擎，需指定监听地址和端口
	dog.Run("127.0.0.1", 8080)
}
```

执行上述代码，请求后应得到响应：

![image-20240217113055414](https://yvling-typora-image-1257337367.cos.ap-nanjing.myqcloud.com/typora/image-20240217113055414.png)





### 模板

PoliteDog提供了简单易用的模板渲染接口。

```go
import (
	"github.com/fangnan700/PoliteDog"
	"net/http"
)

func main() {
	// 创建引擎
	dog := PoliteDog.NewDog()

	// 加载模板
	dog.LoadTemplate("templates/*.html")

	// 创建路由
	router := PoliteDog.NewRouter()
	router.GET("/html", func(ctx *PoliteDog.Context) {
		err := ctx.HTMLTemplate(http.StatusOK, "index.html", "")
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
	})

	// 将路由注册到引擎
	dog.RegisterRouters(router)

	// 启动引擎，需指定监听地址和端口
	dog.Run("127.0.0.1", 8080)
}
```





### 响应数据

#### 1、直接返回

```go
func(ctx *PoliteDog.Context) {
	ctx.Data(http.StatusOK, []byte("data..."))
}
```

#### 2、JSON

```go
func(ctx *PoliteDog.Context) {
	data := map[string]any{
		"name": "admin",
		"age":  21,
	}
	err := ctx.JSON(http.StatusOK, data)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}
```

#### 3、XML

```go
func(ctx *PoliteDog.Context) {
	err := ctx.XML(http.StatusOK, "<tag>Java</tag>")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}
```

#### 4、格式化字符串

```go
func(ctx *PoliteDog.Context) {
	err := ctx.String(http.StatusOK, "%s", "Fuck")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}
```

























