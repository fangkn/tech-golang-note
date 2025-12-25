

# Gin 笔记

 Gin 是使用 golang 实现的 http web 框架。特点： 接口简洁、性能极高。

##  Gin 特性

1、路由快速，不用反射，内在占用, 且可分组。 
2、中间件机制，提高可提扩展性，HTTP请求，可先经过一系列中间件处理，例如：Logger，Authorization，GZIP等。
3、异常处理，服务始终可用，捕获 panic，不会宕机。而且有极为便利的机制处理HTTP请求过程中发生的错误。
4、JSON：Gin可以解析并验证请求的JSON。这个特性对Restful API的开发尤其有用。也支持渲染XML和HTML的渲染。

## 安装

golang 的安装： 

见 [golang 基础](./golang-基础.md) 有安装步骤。

Gin 安装： 

```shell 
go get -u -v github.com/gin-gonic/gin
```

第一个测试程序： hello world 

```go

package main
import "github.com/gin-gonic/gin"

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Gin World !")
	})

	r.Run() 
}
```

说明：

1、`gin.Default()`  生成一个实例，是一个WSGI 应用程序。（Web Server Gateway Interface）
2、`r.GET("/", func(c *gin.Context)` 用于声明了一个路由。什么样的URL 对应什么样的处理函数。并返回什么样的信息。 
3、`r.Run()`  监控端口并把实例跑起来了。如：`r.Run(":1024")`。


## 定义 HTTP 配置

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

http.Server 核心字段有：

- Addr string：监听地址与端口，如 ":8080"。
- Handler Handler：请求处理器。配合 Gin 时设置为 router（*gin.Engine）。
- TLSConfig *tls.Config：自定义 TLS 行为（协议版本、证书、CipherSuite 等）。
- ReadTimeout time.Duration：从连接建立开始，读取整个请求（头+体）的最大时长。流式上传可能不适合设置过小。
- ReadHeaderTimeout time.Duration：仅读取请求头的最大时长，专门防御慢速攻击（slowloris）。
- WriteTimeout time.Duration：写响应的最大时长。若你的响应是长时间流式发送，慎用或留为零值。
- IdleTimeout time.Duration：保持连接（keep-alive）空闲的最大时长。
- MaxHeaderBytes int：请求头最大字节数，默认 1<<20（1MB）。
- ErrorLog *log.Logger：服务器内部错误日志输出目标。
- ConnState func(net.Conn, ConnState)：连接状态变化回调（New/Active/Idle/Hijacked/Closed）。
- BaseContext func(net.Listener) context.Context：为新连接提供基础上下文根，用于统一取消/元数据。
- ConnContext func(ctx context.Context, c net.Conn) context.Context：为每个连接定制上下文（如打标签）。

常用方法：

- ListenAndServe()：在 Addr 上启动 HTTP 服务（明文）。
- ListenAndServeTLS(certFile, keyFile string)：启动 HTTPS，证书文件路径由参数提供。 一般的做法是 https 在 nginx 中配置，然后转发到后方服务，gin 只负责处理业务逻辑。
- Serve(l net.Listener) / ServeTLS(l, certFile, keyFile)：在自定义 net.Listener 上启动服务（可做端口重用、绑定策略等）。
- Shutdown(ctx context.Context)：优雅关闭，拒绝新连接，等待活动请求完成或直到 ctx 超时。
- Close()：立即关闭，活动请求会被打断。
- RegisterOnShutdown(func())：注册关闭时的钩子。
- SetKeepAlivesEnabled(bool)：启用/关闭 HTTP keep-alive。
- RegisterOnShutdown(func())：注册关闭时的钩子。
- SetKeepAlivesEnabled(bool)：启用/关闭 HTTP keep-alive。

 HTTP 方法 有 GET POST PUT PATCH DELETE OPTIONS HEAD 如下：

```go
	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)
```

##  路由 Route 

路由方法有 **GET, POST, PUT, PATCH, DELETE** 和 **OPTIONS**，还有**Any**，可匹配以上任意类型的请求。

### Get 解析路径参数

动态的路由，如 `/player/:name`，通过调用不同的 url 来传入不同的 name。`/player/:name/*role`，`*` 代表可选。

```go
// 匹配 /player/fangkn
r.GET("/player/:name", func(c *gin.Context) {  
	name := c.Param("name")  
	c.String(http.StatusOK, "Hello  %s", name)  
})
```

```sh 
>curl http://127.0.0.1:8080/player/fangkn
Hello  fangkn
```

###  获取Query参数

```go
r.GET("/users", func(c *gin.Context) {
	name := c.Query("name")
	role := c.DefaultQuery("role", "kn")
	c.String(200, "%s is a %s", name, role)
})
```

### 获取POST参数

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
> curl http://127.0.0.1:8080/form -X POST -d 'username=kk&password=1234'
> 	{"password":"1234","username":"kk"}
```

### Query和POST混合参数

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
curl "http://127.0.0.1:8080/posts?id=9876&page=7" -X POST -d 'username=kk&password=1234'

{"id":"9876","page":"7","password":"1234","username":"kk"}
```


### Map参数

这个方式用的比较少。如果不存在会有什么问题？ 

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

## 重定向(Redirect)

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

##  分组路由(Grouping Routes)

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
{"path":"/v1/bag"}
> curl http://127.0.0.1:8080/v2/bag
{"path":"/v2/bag"}
```

##  中间件

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

## BasicAuth 中间件

HTTP Basic Authentication 是一种最简单的身份验证方式。客户端在请求时，通过 HTTP 头部发送用户名和密码 如：

```sh 
Authorization: Basic base64(username:password)
```

```go 
// 模拟一些私人数据
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

// 触发 "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		// 获取用户，它是由 BasicAuth 中间件设置的
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})
```
```sh 
➜  curl -i http://localhost:8080/admin/secrets
HTTP/1.1 401 Unauthorized
Www-Authenticate: Basic realm="Authorization Required"
Date: Mon, 20 Oct 2025 14:25:48 GMT
Content-Length: 0

➜  curl -sS -u austin:1234 http://localhost:8080/admin/secrets
{"secret":{"email":"austin@example.com","phone":"666"},"user":"austin"}

➜ curl -i -H 'Authorization: Basic Zm9vOmJhcg==' http://localhost:8080/admin/secrets
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 20 Oct 2025 14:26:17 GMT
Content-Length: 64

{"secret":{"email":"foo@bar.com","phone":"123433"},"user":"foo"}%

```

见 [gin-basicauth demo](./gin-basicauth/main.go)

## 在中间件中使用 goroutine 

当在中间件或 handler 中启动新的 Goroutine 时，**不能**使用原始的上下文，必须使用只读副本。

```go
func main() {
	r := gin.Default()

	r.GET("/long_async", func(c *gin.Context) {
		// 创建在 goroutine 中使用的副本
		cCp := c.Copy()  // <--------- 必须使用只读副本
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
## 上传文件

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
 见 [gin-upload-file demo](./gin-upload-file/main.go) 

## JSON, JSONP, SecureJSON,PureJSON


** JSONP **

使用 JSONP 向不同域的服务器请求数据。如果查询参数存在回调，则将回调添加到响应体中。是用于跨域请求数据的技术。 	

浏览器出于安全原因（同源策略），默认不允许 JavaScrip发起跨域的 XMLHttpRequest 请求。
但是 `<script>` 标签是个“特例”——它可以加载任意域名的脚本文件。

好处是： 跨域访问数据， 兼容老浏览器 早期不支持 CORS。 

运用的场景有： 

- 跨域获取第三方API数据。 网页上显示天气：https://api.weather.com/data?callback=showWeather	
- 嵌入外部统计或评论系统。 例如早期的百度统计、Disqus 评论都用 JSONP 拉取数据
- 广告系统或新闻聚合。 页面嵌入 `<script>` 加载远程广告、新闻、推荐列表.

替代方案： 现在推荐使用 CORS（跨域资源共享）。 

** SecureJSON  **

使用 SecureJSON 防止 json 劫持。如果给定的结构是数组值，则默认预置 `"while(1),"` 到响应体。

JSON 劫持 是一种前端安全漏洞攻击方式。
攻击者会利用浏览器的特性，通过 `<script>` 签跨域请求接口，从而窃取返回的 JSON 数据。
这在接口直接返回数组的情况下尤其危险！

```html
<script>
    function steal(data) {
        alert("偷到数据: " + data);
    }
</script>
<script src="https://127.0.0.1:8080/users?callback=steal"></script>
```

** PureJSON  **

Gin 默认的 JSON() 会转义 HTML 特殊字符，例如 < 变为 \ u003c。如果要按字面对这些字符进行编码，则可以使用 PureJSON。Go 1.6 及更低版本无法使用此功能。

如： 

```golang
c.JSON(http.StatusOK, gin.H{
	"html": "<b>Hello</b>",
})

// {"html":"\u003cb\u003eHello\u003c/b\u003e"}
```
PureJSON() 与 JSON() 的唯一区别就是：它不会转义 HTML 标签，而是返回原始内容. 

```json 
{"html":"<b>Hello</b>"}

```
代码：

```go

	// 返回 JSON
	router.GET("/json", func(c *gin.Context) {
		data := gin.H{
			"message": "Hello, JSON!",
			"status":  "success",
		}
		c.JSON(http.StatusOK, data)
	})

	// 返回 JSONP
	router.GET("/jsonp", func(c *gin.Context) {
		data := gin.H{
			"message": "Hello, JSONP!",
			"status":  "success",
		}
		// JSONP 需要传 callback 参数，例如：/jsonp?callback=foo
		c.JSONP(http.StatusOK, data)
	})

	// 返回 SecureJSON
	router.GET("/securejson", func(c *gin.Context) {
		data := []string{"Go", "Gin", "Gopher"}
		// SecureJSON 会在 JSON 前加上前缀 "while(1);" 防止 JSON 劫持
		c.SecureJSON(http.StatusOK, data)
	})
```

例子程序 ： [gin-json](./gin-json/main.go)

## jsonp 和 CORS 的跨域

旧版本的浏览器可能不支持CORS ，所以只能用 JSONP 来跨域请求数据。
CORS 的跨域方式

```go
 启用 CORS 支持（允许现代浏览器跨域）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 可改为具体域名，例如 ["https://yourdomain.com"]
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	定义一个统一的接口，自动支持 JSON + JSONP
	r.GET("/data", func(c *gin.Context) {
		data := gin.H{
			"message": "Hello from Gin!",
			"version": "1.0",
			"source":  "Gin JSONP + CORS",
		}

		// 如果有 callback 参数，则返回 JSONP，否则返回 JSON
		callback := c.Query("callback")
		if callback != "" {
			c.JSONP(http.StatusOK, data)
		} else {
			c.JSON(http.StatusOK, data)
		}
	})

```

见 [gin-jsonp-cors](./gin-jsonp-cors/main.go)

前端调用情况：

```js 
fetch("http://localhost:8080/data")
  .then(res => res.json())
  .then(data => console.log("CORS JSON 返回：", data));
```
旧浏览器，使用 `<script>`

```html 
<script>
  function myFunc(data) {
    console.log("JSONP 返回：", data);
  }

  const script = document.createElement("script");
  script.src = "http://localhost:8080/data?callback=myFunc";
  document.body.appendChild(script);
</script>

```
##  html 模板渲染

gin html 模板渲染就是通过加载 html 模板，然后进行变量替代。

加载模板文件 是通过 `LoadHTMLGlob()` 函数。 如下：

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

## 日志

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

## 运行多个服务

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

## 静态文件服务

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


## 优雅地重启或停止

使用 [fvbock/endless](https://github.com/fvbock/endless) 来替换默认的 `ListenAndServe`

见 [golang-endless服务优雅退出.md](./golang-endless服务优雅退出.md)

替代方案:

- [manners](https://github.com/braintree/manners)：可以优雅关机的 Go Http 服务器。
- [graceful](https://github.com/tylerb/graceful)：Graceful 是一个 Go 扩展包，可以优雅地关闭 http.Handler 服务器。
- [grace](https://github.com/facebookgo/grace)：Go 服务器平滑重启和零停机时间部署。


## 热加载调试 Hot Reload

Python 的 `Flask`框架，有 _debug_ 模式，启动时传入 _debug=True_ 就可以热加载(Hot Reload, Live Reload)了。即更改源码，保存后，自动触发更新，浏览器上刷新即可。免去了杀进程、重新启动之苦。

这种热加载在生产环境中基本是要禁止使用的， 生产环境的稳定性，不能允许改完代码直接生效。 

它的一个好处是，在开发环境中， 可以实时看到代码的变化，而不需要重启服务，提高开发效率。

有一个库，可以实现热加载调试。

[air](https://github.com/air-verse/air)：Go 应用的实时热重载工具。

文档说明：[https://github.com/air-verse/air/blob/master/README-zh_cn.md](https://github.com/air-verse/air/blob/master/README-zh_cn.md)

例子见：[hot-reload-air](./hot-reload-air/readme.md)



## 绑定

将 request body 绑定到不同的结构体中 https://gin-gonic.com/zh-cn/docs/examples/bind-body-into-dirrerent-structs/

模型绑定和验证： https://gin-gonic.com/zh-cn/docs/examples/binding-and-validation/

支持 Let's Encrypt https://gin-gonic.com/zh-cn/docs/examples/support-lets-encrypt/


## gin 工程结构

```
your-project/
│
├── cmd/                 # 各种可执行程序入口（main.go 放这里）
│   └── server/
│       └── main.go
│
├── configs/             # 配置文件（yaml/json/toml 等）
│   └── config.yaml
│
├── internal/            # 项目内部代码（不对外暴露）
│   ├── api/             # DTO（请求/响应结构）和路由
│   │   ├── v1/
│   │   │   ├── user.go  # req 和 resp 
│   │   └── router.go    # 注册所有路由
│   │
│   ├── handler/         # controller 层（HTTP 入口）
│   │   └── user_handler.go
│   │
│   ├── service/         # 业务逻辑（service 层）
│   │   └── user_service.go
│   │
│   ├── repository/      # 数据访问（DB / Redis / 外部服务）
│   │   └── user_repo.go
│   │
│   ├── model/           # 领域对象、数据库模型（GORM/SQLX struct）
│   │   └── user.go
│   │
│   ├── middleware/      # Gin 中间件
│   │   ├── auth.go
│   │   └── logger.go
│   │
│   ├── pkg/             # 通用工具库（JWT、加密、日志、雪花ID等）
│   │   ├── jwt/
│   │   ├── logger/
│   │   ├── response/
│   │   └── util/
│   │
│   ├── core/            # 应用核心（启动、初始化、配置、DI）
│   │   ├── config.go
│   │   ├── server.go
│   │   ├── init_db.go
│   │   ├── init_redis.go
│   │   └── init_logger.go
│   │
│   ├── job/             # 定时任务 / 异步任务
│   │   └── clean_user_job.go
│   │
│   └── docs/            # swagger / api 文档
│
├── pkg/                 # 可供外部项目复用的 library（非 internal）
│   └── ...              
│
├── scripts/             # shell 脚本、数据库迁移脚本
│   ├── migrate.sh
│   └── build.sh
│
├── test/                # 集成测试、单元测试
│   └── user_test.go
│
├── go.mod
├── go.sum
└── README.md
```

## 资料

中文学习路径：[https://www.topgoer.com/gin框架/](https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/)

Go Gin 简明教程： [https://geektutu.com/post/quick-go-gin.html](https://geektutu.com/post/quick-go-gin.html)

快速入门：

- [https://geektutu.com/post/quick-go-gin.html](https://geektutu.com/post/quick-go-gin.html)
- [Golang Gin - Github](https://github.com/gin-gonic/gin)
- [https://gin-gonic.com/zh-cn/docs/](https://gin-gonic.com/zh-cn/docs/)

项目学习：

go-admin 有点重了。[https://github.com/go-admin-team/go-admin/blob/master/README.Zh-cn.md](https://github.com/go-admin-team/go-admin/blob/master/README.Zh-cn.md)

Gin 源码分析系列之 Engine 篇：[https://zhuanlan.zhihu.com/p/372097558](https://zhuanlan.zhihu.com/p/372097558)

使用 Gin web 框架的知名项目：

-   [gorush](https://github.com/appleboy/gorush)：Go 编写的通知推送服务器。
-   [fnproject](https://github.com/fnproject/fn)：原生容器，云 serverless 平台。
-   [photoprism](https://github.com/photoprism/photoprism)：由 Go 和 Google TensorFlow 提供支持的个人照片管理工具。
-   [krakend](https://github.com/devopsfaith/krakend)：拥有中间件的超高性能 API 网关。
-   [picfit](https://github.com/thoas/picfit)：Go 编写的图像尺寸调整服务器。
-   [gotify](https://github.com/gotify/server)：使用实时 web socket 做消息收发的简单服务器。
-   [cds](https://github.com/ovh/cds)：企业级持续交付和 DevOps 自动化开源平台。
-   [go-admin](https://github.com/go-admin-team/go-admin): 前后端分离的中后台管理系统脚手架。



