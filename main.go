package main

import (
	// "fmt"
	"log"

	"smartoj-codesandbox/routes"

	"github.com/gin-gonic/gin"
)

// 配置 IP 和 端口号
var HOST, PORT string = "0.0.0.0", "8080"

func createRouter() (*gin.Engine) {
	r := gin.Default()
	// 添加中间件
	r.Use(gin.Recovery()) // 恢复 panic
	return r
}

func main() {
	// 创建 Gin 引擎
	router := createRouter()

	// 创建所有端点
	routes.CreateRoutes(router)

	// 启动服务器
	if err := router.Run(); err != nil {
		log.Fatalf("\033[97;41m Server launched failed! %v \033[0m", err)
	}
}
