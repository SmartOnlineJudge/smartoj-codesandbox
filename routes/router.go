package routes

import (
	"smartoj-codesandbox/routes/judgement"
	"smartoj-codesandbox/routes/health"

	"github.com/gin-gonic/gin"
)

func CreateRoutes(r *gin.Engine) {
	r.GET("/health", health.HealthCheck)  // 健康检查
	r.POST("/sandbox/judgement", judgement.Judge)  // 判题路由
}