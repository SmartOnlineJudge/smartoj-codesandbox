package python

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"smartoj-codesandbox/internal/types"
)


const SOLUTION_FILE_NAME = "solution_code.py"  // 保存用户的解题代码
const MAIN_FILE_NAME = "main.py"  // 主模块
const RUNNER_FILE_NAME = "runner.py"  // 运行器模块
const RUNNER_RELATIVE_PATH = "internal/sandbox/python/" + RUNNER_FILE_NAME  // runner.py 相对路径
const IMPORT_RUNNER_CODE = "from runner import BaseRunner"  // 导入运行器的代码
const RUN_RUNNER_CODE = "Runner().run()"  // 调用运行器的代码
const PYTHON_BIN_NAME = "python3"  // Python 二进制文件名称

// 构建主文件内容 main.py
func constructMainFileContent(judgeTemplate string) string {
	return IMPORT_RUNNER_CODE + "\n" + judgeTemplate + "\n" + RUN_RUNNER_CODE
}

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

	// 构建主文件代码
	mainContent := constructMainFileContent(judgeTemplate)
	fmt.Println(solutionCode)

	solutionFp.WriteString(solutionCode)
	mainFp.WriteString(mainContent)
	
	// 将 runner.py 模块拷贝到指定目录下
	targetRunnerPath := filepath.Join(workspace, RUNNER_FILE_NAME)
	currentDir, _ := os.Getwd()
	currentRunnerPath := filepath.Join(currentDir, RUNNER_RELATIVE_PATH)
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
		limitsInput := fmt.Sprintf("%d %f\n", jd.TimeLimit, jd.MemoryLimit)
		cmd := exec.Command(PYTHON_BIN_NAME, mainPath)
		cmd.Stdin = strings.NewReader(limitsInput + test.InputOutput)
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Status = -2  // -2 为 Runner 内部的错误
			result.Result = string(output)
		} else {
			json.Unmarshal(output, &result)
		}
		result.TestId = test.TestId
		*results = append(*results, result)
	}
	return ""
}
