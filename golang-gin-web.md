
# 简介

 Gin 是使用 golang 实现的 http web 框架。特点： 接口简洁、性能极高。

特性： 

1、 快速。 不用反射，基于 radix，内存战用少。 
2、中间件。 请求会经过一系列中间件处理，提高了框架的可扩展性。
3、异常处理。Gin 可以捕获 panic，并恢复， 服务始终可用，不会宕机。
4、 json 。 Gin 可以解析并验证请求的 json。
5、路径分组。 不同版本的API分组，分组可嵌套，性能不受影响。
6、渲染内置。 原生支持 json xml 和 html 渲染。 

# 安装

**go 的安装**

[[golang-基础]]  有安装步骤。

**Gin 安装**

下载地址：
```shell 
go get -u -v github.com/gin-gonic/gin
```

**hello world** 

第一个测试程序： hello world 

```go
// geektutu.com

// main.go

package main
import "github.com/gin-gonic/gin"

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, xyecho!")
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
```

说明：

1、`gin.Default()`  生成一个实例，是一个WSGI 应用程序。（Web Server Gateway Interface）
2、`r.GET("/", func(c *gin.Context)` 用于声明了一个路由。什么样的URL 对应什么样的处理函数。并返回什么样的信息。 
3、`r.Run()`  监控端口并把实例跑起来了。如：`r.Run(":1024")`。


# 定义 HTTP 配置

定义写超时，读超时等等 如下：

```go
router := gin.Default()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
```

#  HTTP 方法

```go
	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)
```

#  路由 Route 

路由方法有 **GET, POST, PUT, PATCH, DELETE** 和 **OPTIONS**，还有**Any**，可匹配以上任意类型的请求。

## Get 解析路径参数

动态的路由，如 `/player/:name`，通过调用不同的 url 来传入不同的 name。`/player/:name/*role`，`*` 代表可选。

```go
// 匹配 /player/xyecho
r.GET("/player/:name", func(c *gin.Context) {  
	name := c.Param("name")  
	c.String(http.StatusOK, "Hello  %s", name)  
})
```

```sh 
>curl http://127.0.0.1:8080/player/xyecho
Hello  xyecho
```

##  获取Query参数

```go
r.GET("/users", func(c *gin.Context) {
	name := c.Query("name")
	role := c.DefaultQuery("role", "kakaxi")
	c.String(200, "%s is a %s", name, role)
})
```

## 获取POST参数

`PostForm `  从表单中获取参数。
`DefaultPostForm` 从表单中获取参数,如果没有就给出默认的。

```go 
r.POST("/form", func(c *gin.Context) {  
	username := c.PostForm("username")  
	password := c.DefaultPostForm("password", "000000") // 可设置默认值  
  
	c.JSON(http.StatusOK, gin.H{  
		"username": username,  
		"password": password,  
	})  
})
```

测试结果：

```sh 
> curl http://127.0.0.1:8080/form -X POST -d 'username=kakaxi&password=1234'
> 	{"password":"1234","username":"kakaxi"}
```

## Query和POST混合参数

用 `r.POST` 的方式绑定处理函数。 

`GET` 方式的参数用 `c.Query` 和 `c.DefaultQuery` 解析。 
`POST` 方式的参数用   `c.PostForm `和  `c.DefaultPostForm` 解析。 


```go
// GET 和 POST 混合

r.POST("/posts", func(c *gin.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	username := PostForm("username")
	password := c.DefaultPostForm("username", "000000") // 可设置默认值
	
	c.JSON(200, gin.H{
		"id": id,
		"page": page,
		"username": username,
		"password": password,
	})

})
```

测试结果：

```sh 
curl "http://127.0.0.1:8080/posts?id=9876&page=7" -X POST -d 'username=kakaxi&password=1234'

{"id":"9876","page":"7","password":"kakaxi","username":"kakaxi"}
```


## Map参数

这个方式用的比较少。

```go
// Map参数
r.POST("/post2", func(c *gin.Context) {

	ids := c.QueryMap("ids")
	names := c.PostFormMap("names")
	
	c.JSON(http.StatusOK, gin.H{
		"ids": ids,
		"names": names,
	})
})
```

测试结果：

```sh
> curl -g "http://127.0.0.1:8080/post2?ids[kakaxi]=10086&ids[meixi]=10002" -X POST -d 'names[aa]=AAA&names[bb]=BBB'
> {"ids":{"kakaxi":"10086","meixi":"10002"},"names":{"aa":"AAA","bb":"BBB"}}
```

##  Cookie

`c.Cookie() `获取 `cookie` 
`c.SetCookie()` 设置 `cookie` 

```go
import (
    "fmt"

    "github.com/gin-gonic/gin"
)

func main() {

    router := gin.Default()
    router.GET("/cookie", func(c *gin.Context) {
        cookie, err := c.Cookie("gin_cookie")
        if err != nil {
            cookie = "NotSet"
            c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
        }

        fmt.Printf("Cookie value: %s \n", cookie)
    })

    router.Run()
}
```

#  重定向(Redirect)

当我们请求 index1 时，我们可能通过  `Redirect` 重定向到 index2。

用到的函数是  `c.Redirect` 。

```go 
r.GET("/redirect", func(c *gin.Context) {  
    c.Redirect(http.StatusMovedPermanently, "/index")  
})  
  
r.GET("/goindex", func(c *gin.Context) {  
	c.Request.URL.Path = "/"  
	r.HandleContext(c)  
})
```

测试结果如下：

```sh
> curl -i http://127.0.0.1:8080/redirect
HTTP/1.1 301 Moved Permanently
Content-Type: text/html; charset=utf-8
Location: /index
Date: Sun, 05 Mar 2023 04:27:51 GMT
Content-Length: 41

<a href="/index">Moved Permanently</a>.

> curl -i http://127.0.0.1:8080/goindex
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Sun, 05 Mar 2023 04:28:08 GMT
Content-Length: 14

Hello, xyecho!
```

路由重定向，使用 `HandleContext`

```go
r.GET("/test", func(c *gin.Context) {
    c.Request.URL.Path = "/test2"
    r.HandleContext(c)
})
r.GET("/test2", func(c *gin.Context) {
    c.JSON(200, gin.H{"hello": "world"})
})
```


#  分组路由(Grouping Routes)

路由分组功能可以按业务上需求对接口进行分组。在代码和框架上比较清晰。
如：按版本分，部门分，权限分/

```go
// group routes 分组路由
defaultHandler := func(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
	"path": c.FullPath(),
	})
}

// group: v1
v1 := r.Group("/v1")
{
	v1.GET("/gameplayer", defaultHandler)
	v1.GET("/bag", defaultHandler)
}

// group: v2
v2 := r.Group("/v2")
{
	v2.GET("/gameplayer", defaultHandler)
	v2.GET("/bag", defaultHandler)
}
```

测试结果：

```sh 
> curl http://127.0.0.1:8080/v1/gameplayer
{"path":"/v1/gameplayer"}
> curl http://127.0.0.1:8080/v2/gameplayer
{"path":"/v2/gameplayer"}
> curl http://127.0.0.1:8080/v1/bag
{"path":"/v1/bag"}%
> curl http://127.0.0.1:8080/v2/bag
{"path":"/v2/bag"}%
```

#  中间件

```go
func main() {
	// 新建一个没有任何默认中间件的路由
	r := gin.New()

	// 全局中间件
	// Logger 中间件将日志写入 gin.DefaultWriter，即使你将 GIN_MODE 设置为 release。
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery 中间件会 recover 任何 panic。如果有 panic 的话，会写入 500。
	r.Use(gin.Recovery())

	// 你可以为每个路由添加任意数量的中间件。
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// 认证路由组
	// authorized := r.Group("/", AuthRequired())
	// 和使用以下两行代码的效果完全一样:
	authorized := r.Group("/")
	// 路由组中间件! 在此例中，我们在 "authorized" 路由组中使用自定义创建的 
    // AuthRequired() 中间件
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		// 嵌套路由组
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
```

# BasicAuth 中间件

待补充

# 在中间件中使用 goroutine 
当在中间件或 handler 中启动新的 Goroutine 时，**不能**使用原始的上下文，必须使用只读副本。

```go
func main() {
	r := gin.Default()

	r.GET("/long_async", func(c *gin.Context) {
		// 创建在 goroutine 中使用的副本
		cCp := c.Copy()
		go func() {
			// 用 time.Sleep() 模拟一个长任务。
			time.Sleep(5 * time.Second)

			// 请注意您使用的是复制的上下文 "cCp"，这一点很重要
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// 用 time.Sleep() 模拟一个长任务。
		time.Sleep(5 * time.Second)

		// 因为没有使用 goroutine，不需要拷贝上下文
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
```
# 上传文件

上传文件可以分上传单个文件或上传多个文件。 

```go
// 上传文件
// 单个文件上传
r.POST("/upload", func(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	// c.SaveUploadedFile(file, dst)
	c.String(http.StatusOK, "%s uploaded!", file.Filename)
})

// 批量上传
r.POST("/batchupload", func(c *gin.Context) {

	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	
	for _, file := range files {
		log.Println(file.Filename)
		// c.SaveUploadedFile(file, dst)
	}

	c.String(http.StatusOK, "%d files uploaded!", len(files))

})
```


#  JSON

# JSONP

使用 JSONP 向不同域的服务器请求数据。如果查询参数存在回调，则将回调添加到响应体中。

```go
func main() {
	r := gin.Default()

	r.GET("/JSONP", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}
		
		// /JSONP?callback=x
		// 将输出：x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
```

# SecureJSON

使用 SecureJSON 防止 json 劫持。如果给定的结构是数组值，则默认预置 `"while(1),"` 到响应体。是否可以指定其他的？ 

```go
func main() {
	r := gin.Default()

	// 你也可以使用自己的 SecureJSON 前缀
	// r.SecureJsonPrefix(")]}',\n")

	r.GET("/someJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// 将输出：while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
```

# PureJSON

常，JSON 使用 unicode 替换特殊 HTML 字符，例如 < 变为 \ u003c。如果要按字面对这些字符进行编码，则可以使用 PureJSON。Go 1.6 及更低版本无法使用此功能。

```go
func main() {
	r := gin.Default()
	
	// 提供 unicode 实体
	r.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	
	// 提供字面字符
	r.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	
	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
```


#  html 模板渲染

gin html 模板渲染就是通过加载 html 模板，然后进行变量替代。

加载模板文件 是通过 LoadHTMLGlob() 函数。 如下：

```go
package main

import (
"net/http"

"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*")
	r.GET("/index", func(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "火影已经完结了", "name": "我是卡卡西"})
})
	r.Run()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
```

服务启来后，日志上可以看到加载 HTML Templates 的信息。

```sh
[GIN-debug] Loaded HTML Templates (2): 
        - 
        - index.html
```

测试：
```sh
http://127.0.0.1:8080/index
```

![](go-gin-2023-03-05-13-12-12.png)

Gin 默认允许只使用一个 html 模板。 查看[多模板渲染](https://github.com/gin-contrib/multitemplate) 以使用 go 1.6 `block template` 等功能。

# 日志

把输入到屏幕的日志全部打印到日志文件中

```go
// 记录到文件。

f, _ := os.Create("gin.log")

//gin.DefaultWriter = io.MultiWriter(f)

// 如果需要同时将日志写入文件和控制台，请使用以下代码。
gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
```

默认的路由日志格式：

```sh
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /foo                      --> main.main.func1 (3 handlers)
[GIN-debug] GET    /bar                      --> main.main.func2 (3 handlers)
[GIN-debug] GET    /status                   --> main.main.func3 (3 handlers)
```

可以用 `DebugPrintRouteFunc` 来指定日志输出的格式。 


```go
import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r.POST("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, "foo")
	})

	r.GET("/bar", func(c *gin.Context) {
		c.JSON(http.StatusOK, "bar")
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run()
}
```

之后出现的格式是这样的。

```sh 
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

2023/03/05 16:50:41 endpoint POST /foo main.main.func2 3
2023/03/05 16:50:41 endpoint GET /bar main.main.func3 3
2023/03/05 16:50:41 endpoint GET /status main.main.func4 3
```

##  自定义日志文件
```go
func main() {
	router := gin.New()
	// LoggerWithFormatter 中间件会写入日志到 gin.DefaultWriter
	// 默认 gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 你的自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.Run(":8080")
}
```

# 运行多个服务

通过 go 协程可以绑定并启多个服务。

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			},
		)
	})

	return e
}

func router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 02",
			},
		)
	})

	return e
}

func main() {
	server01 := &http.Server{
		Addr:         ":8080",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return server01.ListenAndServe()
	})

	g.Go(func() error {
		return server02.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
```

# 静态文件服务

```go
func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
```


静态资源嵌入使用 [go-assets](https://github.com/jessevdk/go-assets) 将静态资源打包到可执行文件中。
# 资料

中文学习路径：[https://www.topgoer.com/gin框架/](https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/)

快速入门：
- [https://geektutu.com/post/quick-go-gin.html](https://geektutu.com/post/quick-go-gin.html
-   [Golang Gin - Github](https://github.com/gin-gonic/gin)
- [https://gin-gonic.com/zh-cn/docs/](https://gin-gonic.com/zh-cn/docs/)

项目学习 

https://github.com/go-admin-team/go-admin/blob/master/README.Zh-cn.md

Gin 源码分析系列之 Engine 篇：[https://zhuanlan.zhihu.com/p/372097558](https://zhuanlan.zhihu.com/p/372097558)
##### 使用 [Gin](https://github.com/gin-gonic/gin) web 框架的知名项目：

-   [gorush](https://github.com/appleboy/gorush)：Go 编写的通知推送服务器。
-   [fnproject](https://github.com/fnproject/fn)：原生容器，云 serverless 平台。
-   [photoprism](https://github.com/photoprism/photoprism)：由 Go 和 Google TensorFlow 提供支持的个人照片管理工具。
-   [krakend](https://github.com/devopsfaith/krakend)：拥有中间件的超高性能 API 网关。
-   [picfit](https://github.com/thoas/picfit)：Go 编写的图像尺寸调整服务器。
-   [gotify](https://github.com/gotify/server)：使用实时 web socket 做消息收发的简单服务器。
-   [cds](https://github.com/ovh/cds)：企业级持续交付和 DevOps 自动化开源平台。
-   [go-admin](https://github.com/go-admin-team/go-admin): 前后端分离的中后台管理系统脚手架。
# 备忘

## 优雅地重启或停止

使用 [fvbock/endless](https://github.com/fvbock/endless) 来替换默认的 `ListenAndServe`

替代方案:

-   [manners](https://github.com/braintree/manners)：可以优雅关机的 Go Http 服务器。
-   [graceful](https://github.com/tylerb/graceful)：Graceful 是一个 Go 扩展包，可以优雅地关闭 http.Handler 服务器。
-   [grace](https://github.com/facebookgo/grace)：Go 服务器平滑重启和零停机时间部署。


## 热加载调试 Hot Reload

Python 的 `Flask`框架，有 _debug_ 模式，启动时传入 _debug=True_ 就可以热加载(Hot Reload, Live Reload)了。即更改源码，保存后，自动触发更新，浏览器上刷新即可。免去了杀进程、重新启动之苦。

## 绑定

将 request body 绑定到不同的结构体中 https://gin-gonic.com/zh-cn/docs/examples/bind-body-into-dirrerent-structs/

 模型绑定和验证： https://gin-gonic.com/zh-cn/docs/examples/binding-and-validation/

 支持 Let's Encrypt https://gin-gonic.com/zh-cn/docs/examples/support-lets-encrypt/


