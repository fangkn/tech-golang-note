package server

import (
	v1 "gin-demo-service/internal/server/v1"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	stats "github.com/semihalev/gin-stats"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "gin-demo-service/docs"
	"gin-demo-service/internal/middleware"
	"gin-demo-service/internal/version"
)

func NewServer(mode string) *gin.Engine {
	router := gin.New()
	gin.SetMode(mode)

	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.Use(middleware.AccessLog())
	router.Use(middleware.Recovery())
	router.Use(stats.RequestStats())

	router.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, stats.Report())
	})
	pprof.Register(router)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.Info)
	})

	apiv1 := router.Group("/v1/demo")
	apiv1.Use()
	{
		apiv1.GET("/info", v1.GetInfo)
	}
	return router
}
