package judgement

import (
	"fmt"
	"net/http"

	"smartoj-codesandbox/internal/sandbox"
	"smartoj-codesandbox/internal/types"

	"github.com/gin-gonic/gin"
)

// 支持的语言配置
var supportedLanguages = map[string]bool{
	"python": true,
	// "cpp": true, 
	// "c": true,
	// "java": true,
	// "go": true,
}

// 验证判题请求
func validateRequest(jd *types.JudgementData) string {
	if jd.Code == "" {
		return "代码不能为空"
	}
	if jd.Language == "" {
		return "语言类型不能为空"
	}
	if len(jd.Tests) == 0 {
		return "测试用例不能为空"
	}
	if jd.Template == "" {
		return "判题模板不能为空"
	}
	if jd.TimeLimit <= 0 {
		return "时间限制必须大于0"
	}
	if jd.MemoryLimit <= 0 {
		return "内存限制必须大于0"
	}
	if jd.UserId == "" {
		return "用户ID不能为空"
	}

	// 检查语言是否支持
	if _, ok := supportedLanguages[jd.Language]; !ok {
		return fmt.Sprintf("暂不支持 %s 编程语言", jd.Language)
	}

	return ""
}

// 处理评测请求
func Judge(c *gin.Context) {
	var jd types.JudgementData
	
	// 获取请求数据
	err := c.ShouldBindJSON(&jd)
	if err != nil {
		c.JSON(http.StatusOK, types.SandboxResponse{
			Code: 422,
			Message: "无效的请求数据",
		})
		return
	}

	// 验证请求
	errorMessage := validateRequest(&jd)
	if errorMessage != "" {
		c.JSON(http.StatusOK, types.SandboxResponse{
			Code: 422,
			Message: errorMessage,
		})
		return
	}

	// 执行代码
	errorMessage, results := sandbox.ExecuteCode(&jd)
	if errorMessage != "" {
		c.JSON(http.StatusOK, types.SandboxResponse{
			Code: 500,
			Message: errorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, types.SandboxResponse{
		Code: 200,
		Message: "OK",
		Results: results,
	})
}
