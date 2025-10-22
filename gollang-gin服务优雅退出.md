
## 优雅退出

要求是：服务要退出时，要等待所有请求处理完成，再退出。

golang 的服务可以用 endless 库来实现优雅退出。地址：[https://github.com/fvbock/endless](https://github.com/fvbock/endless)

## 运用

可以参考例子： [gin-http-endless](./gin-http-endless/main.go)

```golang 

    // 这个接口会阻塞 sec 秒，用于测试优雅退出
	router.GET("/sleep/:sec", func(c *gin.Context) {
		sec, _ := strconv.Atoi(c.Param("sec"))

		log.Printf("sleep %d seconds", sec)

		time.Sleep(time.Duration(sec) * time.Second)
		c.String(http.StatusOK, "done")
	})

    addr := ":8080"

	// endless.NewServer 返回可热重启的 Server
	srv := endless.NewServer(addr, router)
	// 优先使用实例级超时配置，而不是全局默认值
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 10 * time.Second
	srv.MaxHeaderBytes = 1 << 20

	srv.BeforeBegin = func(add string) {
		log.Printf("Server is starting on %s, pid=%d", add, os.Getpid())
	}

	// 启动服务（支持平滑重启）
	if err := srv.ListenAndServe(); err != nil {
		// endless 在重启/停止时也可能返回错误，这里记录即可
		log.Printf("server error: %v", err)
	}

```

1、用 `endless.NewServer` 定义新的 Server 实例 换掉 `http.Server`

2、写一个阻塞接口，测试优雅退出

3、测试方式：

- 先启动服务
- 用 curl 调用阻塞接口
- 用 `kill -SIGTERM <pid>` 发送信号量
- 确认服务退出时，所有请求都处理完成

如： 

```sh 
./gin-http-endless

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
[GIN-debug] GET    /ping                     --> main.main.func2 (3 handlers)
[GIN-debug] GET    /sleep/:sec               --> main.main.func3 (3 handlers)
2025/10/22 14:59:00 Server is starting on :8080, pid=64112
```
另一个终端：

```sh
curl http://localhost:8080/sleep/50
```

服务终端日志输出：

```sh 
2025/10/22 15:00:11 sleep 50 seconds
```
再启一个终端，更新 gin-http-endless 服务
区别如下： 

```go
	// ping 接口
	//c.JSON(http.StatusOK, gin.H{"message": "pong001"})
	c.JSON(http.StatusOK, gin.H{"message": "pong002"})

	// sleep 接口
	//log.Printf("sleep %d seconds---001", sec)
	//c.String(http.StatusOK, "done-001")
	log.Printf("sleep %d seconds---002", sec)
	c.String(http.StatusOK, "done-002")

```

在 50 秒内 执行以下合集， 重启中 ，服务终端日志输出：

在另一个终端执行 `kill -HUP 64112` ，服务终端日志输出：
```sh 
2025/10/22 15:00:17 64112 Received SIGHUP. forking.
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
[GIN-debug] GET    /ping                     --> main.main.func2 (3 handlers)
[GIN-debug] GET    /sleep/:sec               --> main.main.func3 (3 handlers)
2025/10/22 15:00:18 Server is starting on :8080, pid=64445
2025/10/22 15:00:18 64112 Received SIGTERM.
2025/10/22 15:00:18 64112 Waiting for connections to finish...
2025/10/22 15:00:18 64112 [::]:8080 Listener closed.

```
这里说明服务已经收到信息，要退出，但是在等待所有请求处理完成。`Waiting for connections to finish...`
新的端口已被新服务所用。 

可以请求一下 ping 接口。

```sh 
➜  ~ curl http://localhost:8080/ping
{"message":"pong002"}
```
返回了 pong002 说明服务已更新。旧的服务是 pong001

等 50 秒时间后，服务终端日志输出：

```sh 
2025/10/22 15:01:01 sleep 50 seconds---001
[GIN] 2025/10/22 - 15:01:01 | 200 | 50.001481846s |             ::1 | GET      "/sleep/50"
2025/10/22 15:01:01 64112 Serve() returning...
2025/10/22 15:01:01 server error: accept tcp [::]:8080: use of closed network connection
2025/10/22 15:01:01 server gracefully stopped

```

说明， 旧的服务处理完成sleep 接口的请求了，注意这里的输出是： `sleep 50 seconds---001` 是旧的服务。新是已经更新成 `sleep 50 seconds---002`

这个是执行的流程说了， endless 是可以优雅的执行重启的。不会损坏服务的逻辑的完整性的。 




