package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"smartoj-codesandbox/internal/types"
)

// 必要的文件名称
var SOLUTION_FILE_NAME = "solution_code.py"
var MAIN_FILE_NAME = "main.py"

// 创建相关文件
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
	solutionFp.WriteString(solutionCode)
	mainFp.WriteString(judgeTemplate)

	return ""
}

// 运行 Python 代码
func executePython(workspace string, jd *types.JudgementData, results *types.Results) string {
	errMessage := createCodeFile(workspace, jd.SolutionCode, jd.JudgeTemplate)
	if errMessage != "" {
		return errMessage
	}
	var mainPath = filepath.Join(workspace, MAIN_FILE_NAME)
	for _, test := range jd.Tests {
		cmd := exec.Command("python3", mainPath)
		cmd.Stdin = strings.NewReader(test.InputOutput)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("代码执行异常：\n", string(output))
		} else {
			fmt.Println(string(output))	
		}
	}
	return ""
}
