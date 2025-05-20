package sandbox

import (
	"fmt"
	"path/filepath"

	"smartoj-codesandbox/internal"
)


// 代码执行器
type CodeExecutor struct {
	WorkSpace string // 工作目录
}

// 执行代码
func (e *CodeExecutor) executeCode() {

}

// 创建新的代码执行器
func CreateCodeExecutor(workspace string) (*CodeExecutor) {
	return &CodeExecutor{
		WorkSpace: workspace,
	}
}

// 执行代码
// 提供给外部调用者调用
func ExecuteCode(jd *internal.JudgementData) (error) {
	var workspaceName string = fmt.Sprintf("%d/%s/%s", jd.QuestionId, jd.UserId, jd.Language)
	workspace := filepath.Join("/tmp", workspaceName)
	executor := CreateCodeExecutor(workspace)
	fmt.Println(executor.WorkSpace)
	return nil
}
