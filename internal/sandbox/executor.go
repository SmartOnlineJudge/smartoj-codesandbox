package sandbox

import (
	"os"

	"smartoj-codesandbox/internal/types"
)

// 定义代码执行函数映射
var executeFuncMapping = map[string]func(string, *types.JudgementData, *types.Results) string{
	"python": executePython,
}

// 代码执行器
type CodeExecutor struct {
	workspace string  // 工作目录
	jd *types.JudgementData  // 判题元数据
}

// 执行代码
func (e *CodeExecutor) execute() (string, *types.Results) {
	// 创建工作目录
	err := os.MkdirAll(e.workspace, 0755)

	// 创建判题结果
	results := make(types.Results, 0, len(e.jd.Tests))

	if err != nil {
		return "目录创建失败", &results
	}

	// 获取判题函数
	executeFunc, ok := executeFuncMapping[e.jd.Language]
	if !ok {
		return "选择了不支持的编程语言", &results
	}

	// 执行判题
	errMessage := executeFunc(e.workspace, e.jd, &results)
	if errMessage != "" {
		return errMessage, &results
	}
	return "", &results
}
