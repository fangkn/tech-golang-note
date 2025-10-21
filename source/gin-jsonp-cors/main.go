package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 可改为具体域名，例如 ["https://yourdomain.com"]
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/data", func(c *gin.Context) {
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
