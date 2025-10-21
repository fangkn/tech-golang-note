//go:build !windows
// +build !windows

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	// 生产环境可使用 ReleaseMode（减少日志）
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 路由
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin World!")
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
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

	log.Println("server gracefully stopped")
}
