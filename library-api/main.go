package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/example/library-api/config"
	"github.com/example/library-api/database"
	"github.com/example/library-api/middleware"
	"github.com/example/library-api/routes"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	database.InitDB()

	// 初始化限流中间件
	middleware.InitRateLimiter(cfg)

	// 创建Gin引擎
	router := gin.Default()

	// 注册中间件
	router.Use(middleware.RateLimitMiddleware())

	// 注册路由
	routes.RegisterRoutes(router)

	// 打印所有注册的路由
for _, route := range router.Routes() {
	log.Printf("已注册路由: %s %s\n", route.Method, route.Path)
}

// 启动服务器
log.Printf("服务器启动在端口 %s", cfg.ServerPort)
if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
	log.Fatalf("服务器启动失败: %v", err)
}
}