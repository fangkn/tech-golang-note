# 多版本问题

遇到一个问题如下：

```sh 
# encoding
compile: version "go1.18.2" does not match go tool version "go1.18.2"
# internal/goarch
compile: version "go1.20.1" does not match go tool version "go1.18.2"
```

原因是系统的go 版本是 `go1.18.2` 而第三方库要求的是 `go1.20.1`。  

系统安装了多个版本的 go 。而被使用是 `go1.18.2`。由于这个版本其他同事可能在使用，

所以想在不影响他人使用 `go1.18.2`的情况下，在我自己的用户下使用 `go1.20.1`

试着把 `go1.20.1` 的路径配置到用户下的 `.bashrc` 或 `.profile`。

直接执行 `go verison` 可使用 `go1.20.1` 。但是在makefile 时，却还是使用了 `go1.18.2`的版本。

猜想原因是：makefile 起子shell 环境变量 和 我所在用户环境变量不一样。

解决方法： 在当前的用户下建立 bin 把   `go1.20.1` 连接到 bin 下:` ln -s /usr/local/go/bin/go go `

再次尝试，可以使用。

# beego dev note 2023-04-17

测试时，发现 beego 读取配置文件时，不会实时加载配置文件 ，而是在启动服务时一次加载后，放到 map 中去，然后 beego.AppConfig.String 时再从 map 中取出来。 

from :https://yushuangqi.com/blog/2016/beego-webframework-config.html

一个热加载的问题：https://stackoverflow.com/questions/36078560/how-beegogo-app-framework-will-reload-the-application-if-there-is-any-change-i

https://blog.csdn.net/wangzhufei/article/details/93144116

查了一下资料， beego 好像无法做热加载配置文件。

beego 框架错误时 panic 返回一些 html 文本。 有些不友好。 

from : [https://www.jianshu.com/p/1eab4b981559](https://www.jianshu.com/p/1eab4b981559)



# beego v2  读配置时遇到问题

用 beego v2 想测试一下，实时加载配置的方式： Port := beego.AppConfig.String("mysql_port")

结果却遇到了这个问题。也看不出什么来。

```sh
c:\4K\github\golang-learning-pro\beego-chat>go run main.go
C:\Users\14476\AppData\Local\Temp\go-build2617687394\b001\exe\main.exe flag redefined: graceful
panic: C:\Users\14476\AppData\Local\Temp\go-build2617687394\b001\exe\main.exe flag redefined: graceful
goroutine 1 [running]:
flag.(*FlagSet).Var(0xc000068180, {0x10ebee8, 0x156942f}, {0xff3c77, 0x8}, {0x1009fb7, 0x21})
​    C:/Program Files/Go/src/flag/flag.go:982 +0x2f9
flag.BoolVar(...)
​    C:/Program Files/Go/src/flag/flag.go:723
[github.com/beego/beego/v2/server/web/grace.init.0](http://github.com/beego/beego/v2/server/web/grace.init.0)()
​    C:/Users/14476/go/pkg/mod/[github.com/beego/beego/v2@v2.0.7/server/web/grace/grace.go:93](http://github.com/beego/beego/v2@v2.0.7/server/web/grace/grace.go:93) +0x52

exit status 2
```


看了一下beego v2 读取配置的方式，是：

```golang
import   "[github.com/beego/beego/v2/core/config](http://github.com/beego/beego/v2/core/config)"

dbMap, _ := config.GetSection("mysql")

Port := dbMap["mysql_port"]
```

改成这种方式后就可以了。 

这个变化还真有点大。

# 对于中文等字符数的判断 

```golang
if len([]rune(req.Name)) > 100
```


在 Go 语言中， 这段代码的含义是检查 `req.Remark` 字符串转换成 rune 切片后的长度是否超过 100。如果超过 100，那么这个条件判断就会返回 `true`。

这里的关键点是 `[]rune(req.Name)`。`rune` 在 Go 语言中是一个别名，代表 UTF-8 编码的 Unicode 码点。在 Go 中，字符串是由字节组成的，而不是字符。当处理包含多字节字符（如中文、表情符号等）的字符串时，直接使用 `len()` 函数会返回字符串的字节数量，而不是字符数量。为了正确地计算这些字符串的字符数量，我们需要将其转换为 rune 切片，然后再使用 `len()` 函数。

例如，如果 `req.Name` 包含中文字符，那么每个中文字符通常会占用 3 个字节。但是，当我们将其转换为 rune 切片后，每个中文字符都会被视为一个元素，因此 `len([]rune(req.Name))` 会返回正确的字符数量。



# 二制进流如何在结构体中定义

接入小程序创建二维码时，微信平台返回了一个二维码的二进制流图片。应该如何处理比较合适。

>注意微信返回的数据，官方文档写的是 `Buffer`， 我们在接收的时候，判断只要不是 `json` 格式，即获取到 `header` 并且获取到数据流，将数据流做 `base64_encode`，然后拼接在`data:image/jpeg;base64,` 后面即构成了图片的完整 `base64` 编码。 其中 `image/jpeg` 是我们从 `header` 中拿到的。

> 通过该种方式，我们即不需要在服务器上存储任何图片信息，以造成负担。当然如果有特殊需求的另当别论。

# 编译出现 ： import cycle not allowe

```sh
package cag-wx-proxy
        imports cag-wx-proxy/internal/config
        imports cag-wx-proxy/internal/config: import cycle not allowed
package cag-wx-proxy
        imports cag-wx-proxy/internal/config
        imports cag-wx-proxy/internal/config: import cycle not allowed
        
# 这个问题是在一个模块引用自己的 package 
#
```

在Golang中，当你遇到“import cycle not allowed”的错误时，这意味着存在循环导入问题，即两个或多个包之间相互依赖，形成了一个闭环。这违反了Golang的设计原则，因为Golang程序必须是无环的。

可以有 godepgraph 工具 查出依赖关系 

[https://pkg.go.dev/github.com/alovn/godepgraph#section-readme](https://pkg.go.dev/github.com/alovn/godepgraph#section-readme)

**golang pb 协议如何定义返回错误码和错误信息** 

[https://jiajunhuang.com/articles/2022_11_04-grpc_error_handling.md.html](https://jiajunhuang.com/articles/2022_11_04-grpc_error_handling.md.html)


# window go mod  权限问题

安装后，配置在  c:\Program Files 下， go mod tidy  在下载库时，用户没有权限会报错误如下：

  ```cmd
go: writing go.mod cache: mkdir c:\Program Files\Go\pkg\mod\cache: Access is denied.
go: writing go.mod cache: mkdir c:\Program Files\Go\pkg\mod\cache: Access is denied.

  ```

对  c:\Program Files\Go 添加当前用户的读写和执行权限

# 对 map 的 value 进行类型转换


```go
	Config1 := mapConf["game"].(map[string]interface{}) // 有问题
	for key, value := range confMap {
		Config1[key] = value
	}
```
mapConf["game"]为nil：如果mapConf["game"]是 nil，那么在尝试进行类型断言时也会导致 panic。

改成：

```go
    Config1 := make(map[string]interface{})
	for key, value := range confMap {
		Config1[key] = value
	}
	mapConf["game"] = Config1
```

**错误：invalid version: git ls-remote -q origin in**

```
internal\logic\checkPhoneCaptchaLogic.go:16:2: gitlab.vrviu.com/cloud_esport_backend/saas_user_server@v0.0.0-20240524031850-03bf96fbcf45: invalid version: git ls-remote -q origin in C:\Users\kane\go\pkg\mod\cache\vcs\b6d3346b9a08a643b6538f20717fbbcf2a9651f4d8a8922abe6346170ec83489: exit status 128:
        warning: missing OAuth configuration for gitlab.vrviu.com - see https://aka.ms/gcm/gitlab for more information
        fatal: Cannot prompt because user interactivity has been disabled.
        fatal: could not read Username for 'https://gitlab.vrviu.com': terminal prompts disabled
Confirm the import path was entered correctly.
```

![](./assets/go-2024-10-02_16-11-34.jpg)

在 git 的全局配置中加下面的配置 
```
[url "https://github"]
	insteadOf = git://github
[url "git@gitlab.vrviu.com:"]
        insteadOf = https://gitlab.vrviu.com/   
```