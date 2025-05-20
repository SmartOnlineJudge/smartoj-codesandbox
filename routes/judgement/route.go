package judgement

import (
	"net/http"

	"smartoj-codesandbox/internal"
	// "smartoj-codesandbox/internal/sandbox"

	"github.com/gin-gonic/gin"
)

// 验证判题请求
func validateRequest(jd *internal.JudgementData) string {
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
	if _, ok := internal.SupportedLanguages[jd.Language]; !ok {
		return "不支持的语言类型: " + jd.Language
	}

	return ""
}

// 处理评测请求
func Judge(c *gin.Context) {
	var jd internal.JudgementData
	
	// 获取请求数据
	err := c.ShouldBindJSON(&jd)
	if err != nil {
		c.JSON(http.StatusOK, internal.SandboxResponse{
			Code: 422,
			Message: "无效的请求数据",
		})
		return
	}

	// 验证请求
	errorMessage := validateRequest(&jd)
	if errorMessage != "" {
		c.JSON(http.StatusOK, internal.SandboxResponse{
			Code: 422,
			Message: errorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, internal.SandboxResponse{
		Code: 200,
		Message: "OK",
	})
	
	// 执行代码
	// sandbox.ExecuteCode(&jd)
}
