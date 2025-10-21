package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Gin 路由
	router := gin.Default()

	// 简单路由
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin World!")
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

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

	// 提供字面字符
	router.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	// 自定义 http.Server
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
			panic(err)
		}
	}()

	// 捕获 SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	// 可选：在关闭时做清理
	s.RegisterOnShutdown(func() {
		log.Println("shutting down: cleanup tasks...")
	})

	// 设定最多 5 秒完成在途请求
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("server gracefully stopped")
}
