package sandbox

import (
	"fmt"
	"path/filepath"
	"os"

	"smartoj-codesandbox/internal/types"
)


// 创建代码执行器
func CreateCodeExecutor(workspace string, jd *types.JudgementData) (*CodeExecutor) {
	return &CodeExecutor{
		workspace: workspace,
		jd: jd,
	}
}

// 执行代码
// 提供给外部调用者调用
func ExecuteCode(jd *types.JudgementData) (string, *types.Results) {
	var workspaceName string = fmt.Sprintf("%d/%s/%s", jd.QuestionId, jd.UserId, jd.Language)
	workspace := filepath.Join(os.TempDir(), workspaceName)
	executor := CreateCodeExecutor(workspace, jd)
	return executor.execute()
}
