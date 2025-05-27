package python

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"smartoj-codesandbox/internal/types"
)

// 必要的文件名称
var SOLUTION_FILE_NAME = "solution_code.py"
var MAIN_FILE_NAME = "main.py"
var RUNNER_FILE_NAME = "runner.py"

// Python 二进制文件名称
var PYTHON_BIN_NAME = "python3"

// 创建必要的代码文件
func createCodeFile(workspace, solutionCode, judgeTemplate string) string {
	// 构建文件路径
	var solutionPath, mainPath string
	solutionPath = filepath.Join(workspace, SOLUTION_FILE_NAME)
	mainPath = filepath.Join(workspace, MAIN_FILE_NAME)

	// 获取文件句柄
	solutionFp, err := os.Create(solutionPath)
	if err != nil {
		return "解题函数文件创建失败"
	}
	mainFp, err := os.Create(mainPath)
	if err != nil {
		return "判题模板文件创建失败"
	}

	// 关闭文件句柄
	defer func() {
		solutionFp.Close()
		mainFp.Close()
	}()

	// 创建相关文件
	judgeTemplate += "\n_runner = Runner(solution)\n_runner.run()"  // 增加判题模板的调用
	solutionFp.WriteString(solutionCode)
	mainFp.WriteString(judgeTemplate)
	targetRunnerPath := filepath.Join(workspace, RUNNER_FILE_NAME)
	currentDir, _ := os.Getwd()
	currentRunnerPath := filepath.Join(currentDir, "internal/sandbox/python/runner.py")
	exec.Command("cp", currentRunnerPath, targetRunnerPath).Run()

	return ""
}

// 运行 Python 代码
func ExecutePython(workspace string, jd *types.JudgementData, results *types.Results) string {
	errMessage := createCodeFile(workspace, jd.SolutionCode, jd.JudgeTemplate)
	if errMessage != "" {
		return errMessage
	}
	var mainPath = filepath.Join(workspace, MAIN_FILE_NAME)
	for _, test := range jd.Tests {
		var result types.Result
		cmd := exec.Command(PYTHON_BIN_NAME, mainPath)
		cmd.Stdin = strings.NewReader(test.InputOutput)
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Status = -1
			result.Result = string(output)
		} else {
			json.Unmarshal(output, &result)
		}
		result.TestId = test.TestId
		*results = append(*results, result)
	}
	return ""
}
