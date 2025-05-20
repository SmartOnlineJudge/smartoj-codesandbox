package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"code": 200,
	})
}
