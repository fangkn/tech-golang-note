package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/mcp"
)

func main() {
	// 加载配置
	var c mcp.McpConf
	conf.MustLoad("config.yaml", &c)

	// 创建 MCP 服务器
	server := mcp.NewMcpServer(c)
	defer server.Stop()

	// 注册计算器工具
	calculatorTool := mcp.Tool{
		Name:        "calculator",
		Description: "执行基础数学运算",
		InputSchema: mcp.InputSchema{
			Properties: map[string]any{
				"operation": map[string]any{
					"type":        "string",
					"description": "要执行的操作 (add, subtract, multiply, divide)",
					"enum":        []string{"add", "subtract", "multiply", "divide"},
				},
				"a": map[string]any{
					"type":        "number",
					"description": "第一个操作数",
				},
				"b": map[string]any{
					"type":        "number",
					"description": "第二个操作数",
				},
			},
			Required: []string{"operation", "a", "b"},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			var req struct {
				Operation string  `json:"operation"`
				A         float64 `json:"a"`
				B         float64 `json:"b"`
			}

			if err := mcp.ParseArguments(params, &req); err != nil {
				return nil, fmt.Errorf("参数解析失败: %v", err)
			}

			// 执行操作
			var result float64
			switch req.Operation {
			case "add":
				result = req.A + req.B
			case "subtract":
				result = req.A - req.B
			case "multiply":
				result = req.A * req.B
			case "divide":
				if req.B == 0 {
					return nil, fmt.Errorf("除数不能为零")
				}
				result = req.A / req.B
			default:
				return nil, fmt.Errorf("未知操作: %s", req.Operation)
			}

			// 返回格式化结果
			return map[string]any{
				"expression": fmt.Sprintf("%g %s %g", req.A, getOperationSymbol(req.Operation), req.B),
				"result":     result,
			}, nil
		},
	}

	// 注册工具到服务器
	if err := server.RegisterTool(calculatorTool); err != nil {
		log.Fatalf("注册计算器工具失败: %v", err)
	}

	fmt.Printf("启动 MCP 服务器，端口: %d\n", c.Port)
	server.Start()
}

func getOperationSymbol(op string) string {
	switch op {
	case "add":
		return "+"
	case "subtract":
		return "-"
	case "multiply":
		return "×"
	case "divide":
		return "÷"
	default:
		return op
	}
}
