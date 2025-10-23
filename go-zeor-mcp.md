# golang go-zore mcp 

mcp 的原理，是通过 http 协议，将请求发送到 mcp server，mcp server 会根据请求，调用对应的工具，工具的执行结果，会返回给 mcp server，mcp server 会将结果返回给请求方。

详细见其他文档。

这个是 go-zore 框架的 mcp 示例。一个很简单的计算器。 

文档见： [Golang 快速开发 MCP Servers](https://zhuanlan.zhihu.com/p/1904921186875450778)

代码：[go-zore-mcp-demo](https://github.com/fangkn/tech-golang-note/tree/main/source/go-zore-mcp-demo)

编译完成之后，运行 `./calculator-assistant` 即可。

```sh
./calculator-assistant

启动 MCP 服务器，端口: 8080

```

在 trae 配置 mcp server。 

设置图标 -> MCP-> 添加 -> 手动添加

![](./assets/golang-mcp-2025-10-23_17-06-10.png)

配置如下： 

```json
{
  "mcpServers": {
    "calculator": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "http://localhost:8080/sse"
      ]
    }
  }
}
```

确认之后， 看是否已经有绿色的打勾图标。 

在 trae 尝试使用 mcp 计算器。

![](./assets/golang-mcp-2025-10-23_16-50-40.png)

mcp 的服务 calculator-assistant 输出的日志如下： 

 ```sh 
{"@timestamp":"2025-10-23T16:50:40.982+08:00","caller":"handler/loghandler.go:167","content":"[HTTP] 202 - POST /message?session_id=dd75c5ad-e500-4cd9-a2f8-1762fa23259e - [::1]:60645 - node","duration":"0.5ms","level":"info","span":"bf2baad4fcdbc863","trace":"59afe2a27aa9bca6de9a14e072319951"}
{"@timestamp":"2025-10-23T16:50:40.989+08:00","caller":"mcp/server.go:212","content":"Received tools call request with ID: 3","level":"info"}
{"@timestamp":"2025-10-23T16:50:40.990+08:00","caller":"mcp/server.go:631","content":"Executing tool 'calculator' with arguments: map[string]interface {}{\"a\":100, \"b\":200, \"operation\":\"add\"}","level":"info"}
{"@timestamp":"2025-10-23T16:50:40.991+08:00","caller":"mcp/server.go:762","content":"Tool call result: mcp.CallToolResult{Result:mcp.Result{Meta:map[string]interface {}(nil)}, Content:[]interface {}{mcp.typedTextContent{Type:\"text\", TextContent:mcp.TextContent{Text:\"{\\\"expression\\\":\\\"100 + 200\\\",\\\"result\\\":300}\", Annotations:(*mcp.Annotations)(0xc00041e320)}}}, IsError:false}","level":"info"}
 ```

 从日志中可以看出， 有 `{\"a\":100, \"b\":200, \"operation\":\"add\"}` 计算器成功计算出了 100 + 200 = 300

