package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
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

	// 注册 IP 定位工具
	ipLocationTool := mcp.Tool{
		Name:        "ip-location",
		Description: "查询 IP 地址的位置信息",
		InputSchema: mcp.InputSchema{
			Properties: map[string]any{
				"ip": map[string]any{
					"type":        "string",
					"description": "要查询的 IP 地址",
				},
			},
			Required: []string{"ip"},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			var req struct {
				IP string `json:"ip"`
			}
			if err := mcp.ParseArguments(params, &req); err != nil {
				return nil, fmt.Errorf("参数解析失败: %v", err)
			}

			// 打开 SQLite 数据库
			db, err := sql.Open("sqlite3", "./ip_database.db")
			if err != nil {
				return nil, fmt.Errorf("打开数据库失败: %v", err)
			}
			defer db.Close()

			// 将 IPv4 字符串转为 uint32
			ipInt, err := ipv4ToUint32(req.IP)
			if err != nil {
				return nil, fmt.Errorf("IP 格式不正确: %v", err)
			}

			// 优先尝试使用数值区间字段查询
			var (
				country  string
				province string
				city     string
				district string
				isp      string
			)
			query := `SELECT country, province, city, district, isp FROM ip_data WHERE ? BETWEEN start_ip_int AND end_ip_int LIMIT 1`
			err = db.QueryRowContext(ctx, query, ipInt).Scan(&country, &province, &city, &district, &isp)
			if err != nil {
				// 如果没有数值字段，尝试使用 id 回退策略
				query2 := `SELECT country, province, city, district, isp FROM ip_data WHERE id <= ? ORDER BY id DESC LIMIT 1`
				if scanErr := db.QueryRowContext(ctx, query2, ipInt).Scan(&country, &province, &city, &district, &isp); scanErr != nil {
					// 最后回退：全表扫描并在代码中比较 start_ip/end_ip 范围
					rows, qerr := db.QueryContext(ctx, `SELECT start_ip, end_ip, country, province, city, district, isp FROM ip_data`)
					if qerr != nil {
						return nil, fmt.Errorf("查询失败: %v", qerr)
					}
					defer rows.Close()

					var found bool
					for rows.Next() {
						var startIP, endIP string
						if err := rows.Scan(&startIP, &endIP, &country, &province, &city, &district, &isp); err != nil {
							return nil, fmt.Errorf("扫描数据失败: %v", err)
						}
						startInt, serr := ipv4ToUint32(startIP)
						endInt, eerr := ipv4ToUint32(endIP)
						if serr == nil && eerr == nil && ipInt >= startInt && ipInt <= endInt {
							found = true
							break
						}
					}
					if !found {
						return nil, fmt.Errorf("未找到匹配记录")
					}
				}
			}

			return map[string]any{
				"ip":       req.IP,
				"country":  country,
				"province": province,
				"city":     city,
				"district": district,
				"isp":      isp,
			}, nil
		},
	}

	// 注册工具到服务器
	if err := server.RegisterTool(ipLocationTool); err != nil {
		log.Fatalf("注册 IP 定位工具失败: %v", err)
	}

	fmt.Printf("启动 MCP 服务器，端口: %d\n", c.Port)
	server.Start()
}

func ipv4ToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP: %s", ipStr)
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}
	return uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3]), nil
}
