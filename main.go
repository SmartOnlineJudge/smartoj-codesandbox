package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func GetDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func main() {
	cli, err := GetDockerClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	ctx := context.Background()

	// 1. 定义两个 Python 文件的内容（Go 变量）
	solutionCode := `
def calculate():
    return 100 / 2`

	mainCode := `
from solution import calculate
n = int(input())
nums = []
for _ in range(n):
    nums.append(int(input()))
print(f"nums is: {nums}")
print(calculate())`

	inputData := "3\\n1\\n2\\n3"

	// 2. 定义容器命令：生成文件并运行 main.py
	cmd := []string{
		"sh", "-c",
		fmt.Sprintf(
			"mkdir -p /workspace && echo '%s' > /workspace/solution.py && echo '%s' > /workspace/main.py && echo -e '%s' | python3 /workspace/main.py",
			solutionCode, // 写入 solution.py
			mainCode,     // 写入 main.py
			inputData,
		),
	}

	imageName := "python:3.11-alpine"

	config := container.Config{
		Image: imageName,
		Cmd: cmd,
		Tty: false, // 确保禁用TTY
	}

	resp, err := cli.ContainerCreate(ctx, &config, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	// 等待容器执行完毕
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	// 获取容器日志
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,      // 确保获取完整日志
		Details:    false,
	})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	var stderrBuf, stdoutBuf bytes.Buffer

	// 使用标准库的 stdcopy 处理日志输出
	_, err = stdcopy.StdCopy(&stderrBuf, &stdoutBuf, out)
	if err != nil {
		panic(err)
	}

	stdoutStr := stdoutBuf.String()
	stderrStr := stderrBuf.String()

	// 删除容器
	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{}); err != nil {
		panic(err)
	}
	fmt.Printf("标准输出：\n")
	fmt.Println(stderrStr)
	fmt.Printf("程序报错：\n")
	fmt.Println(stdoutStr)
}
