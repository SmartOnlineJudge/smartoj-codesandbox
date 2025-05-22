# 智能算法刷题平台——代码沙箱仓库

## 判题沙箱介绍

判题沙箱指的是专门运行给定代码的一个应用程序。外部程序需要使用 HTTP 协议与判题沙箱进行通信。

它负责接收用户编写的函数、判题模板、输入输出数据、编程语言类型、内存限制、时间限制等等参数，最终判题沙箱进程通过一个子协程来对本次输入进行判题操作。

沙箱进程本身部署在 Docker 容器中，具有天生的隔离性，保证了服务器的安全性。

一句话总结判题沙箱的作用：运行用户提交的代码，判断正误并将结果返回。

## 技术栈
1. Goroutine（Go 语言轻量级并发协程）
2. Gin（Go 语言轻量级兼高性能 Web 后端框架）
3. Docker（容器隔离技术）

## 沙箱请求与响应
### 沙箱调用过程
![](https://cdn.nlark.com/yuque/0/2025/png/47866636/1747230481656-709115f0-20f2-484f-89ab-508f39469b83.png)

### 沙箱内部调用过程
![](https://cdn.nlark.com/yuque/0/2025/png/47866636/1747280148859-9951a365-9d64-4915-833f-4d4cf0c8a2ab.png)

## 沙箱设计思路
### 输入输出数据
为了适应所有的编程语言，所有的输入和输出均保存在同一个字符串中，并且使用最原始的方式保存数据。

例如真实的输入输出：

```plain
3
1 2 3
6
```

实际保存的输入输出：

```plain
"3\n1 2 3\n6"
```

这样的方式可以适应所有的编程语言，因为几乎所有的编程语言都支持从标准输入流中读取数据，当有了这些数据之后，相应的编程语言可以根据它的特征进一步构造好数据类型，再将构造好的输入数据传递给用户编写的解题函数，然后捕获这个函数的输出，最后再用从输入流中读取的正确输出与该函数的输出进行比对。最后由代码沙箱来捕获判题结果并返回给调用者。

### 运行代码
当拿到判题模板和解题函数以后，需要将这两份代码前保存在`/tmp`目录下，这样就会比较方便运行代码。

除此之外，调用者还会将`题目ID`和`用户ID`传过来，因此可以使用这两个值来处理目录冲突的问题。

期望的目录格式：`/tmp/<question_id>/<user_id>/<language>`。

假设接收到数据是这样的：

```json
{
  "language": "python",
  "user_id": "2025012849102",
  "question_id": 1,
  ...
}
```

那么这时应该先创建一个目录：`/tmp/1/2025012849102/python`。

然后在这个目录的内部创建相关的代码文件。

例如：

+ `/tmp/1/2025012849102/python/solution_code.py`
+ `/tmp/1/2025012849102/python/main.py`

最后由 Go 语言调用 Python 解释器并将输入输出数据传入主函数中。

```bash
echo -e '<输入输出文本>' | python3 /tmp/1/2025012849102/python/main.py
```

### 运行结果与异常捕获
对于所有的判题模板，应该使用相同格式的输出（JSON字符串）：

```json
{
  "status": 1,  // 代码运行异常时为-1
  "result": "OK | ErrorMessage",
  "is_success": true,
  "solution": "解题函数的输出（字符串）",
  "criterion": "正确的输出（字符串）",
  "time_consumed": "运行该测试用例消耗的时间",
  "memory_consumed": "运行该测试用例消耗的内存",
}
```

最后由 Go 语言捕获判题模板的输出，并将其加入本次请求对应的响应中。

### 限制代码运行内存和时间
这个任务由判题模板来限制是最合适的，并且限制粒度需要控制到函数级别。

也就是说需要对`solution`进行单独限制，由相应的编程语言来限制。

假设用户编写的`solution`代码是：

```python
def solution():
    while True:
        pass  # 模拟长时间/高内存操作
```

以下是各个编程语言的内存时间限制示例：

#### Python（只支持 Linux 系统）
```python
import concurrent.futures
import resource
import time
import os

def solution():
    while True:
        time.sleep(0.01)

def limited_solution():
    # 设置内存限制（10MB）
    soft, hard = 10 * 1024 * 1024, 15 * 1024 * 1024
    resource.setrlimit(resource.RLIMIT_AS, (soft, hard))
    solution()

def main():
    with concurrent.futures.ProcessPoolExecutor() as executor:
        future = executor.submit(limited_solution)
        try:
            future.result(timeout=3)
        except concurrent.futures.TimeoutError:
            print("Timeout!")
            future.cancel()
        except Exception as e:
            print("Error:", e)

if __name__ == "__main__":
    main()
```

#### C
```c
#include <stdio.h>
#include <signal.h>
#include <unistd.h>
#include <sys/resource.h>

void handle_timeout(int sig) {
    printf("Timeout!\n");
    _exit(1);
}

void solution() {
    while (1) { /* 模拟死循环 */ }
}

int main() {
    // 设置内存限制（10MB）
    struct rlimit rl;
    rl.rlim_cur = 10 * 1024 * 1024; // 10MB
    rl.rlim_max = 10 * 1024 * 1024;
    setrlimit(RLIMIT_AS, &rl);

    // 设置时间限制（3秒）
    signal(SIGALRM, handle_timeout);
    alarm(3);

    solution(); // 执行目标函数
    return 0;
}
```

#### C++
```cpp
#include <iostream>
#include <future>
#include <chrono>
#include <thread>
#include <sys/resource.h>

int solution() {
    while (true) std::this_thread::sleep_for(std::chrono::milliseconds(10));
    return 42;
}

int main() {
    // 设置内存限制（10MB）
    struct rlimit rl;
    rl.rlim_cur = 10 * 1024 * 1024;
    rl.rlim_max = 10 * 1024 * 1024;
    setrlimit(RLIMIT_AS, &rl);

    // 设置时间限制（3秒）
    std::future<int> result = std::async(std::launch::async, solution);
    auto status = result.wait_for(std::chrono::seconds(3));

    if (status == std::future_status::ready)
        std::cout << "Result: " << result.get() << std::endl;
    else
        std::cerr << "Function timed out!" << std::endl;

    return 0;
}
```

#### Java
```java
import java.util.concurrent.*;

public class Main {
    public static int solution() {
        while (true); // 模拟无限循环
        // return 42;
    }

    public static void main(String[] args) throws Exception {
        ExecutorService executor = Executors.newSingleThreadExecutor();
        Future<Integer> future = executor.submit(() -> solution());

        try {
            System.out.println(future.get(3, TimeUnit.SECONDS));
        } catch (TimeoutException e) {
            System.out.println("Timed out!");
            future.cancel(true);
        } finally {
            executor.shutdownNow();
        }
    }
}
```

#### Go（编程语言版本需要 >= 1.19）
```go
package main

import (
    "context"
    "fmt"
    "runtime/debug"
    "time"
)

func init() {
    // 设置内存限制为 10MB
    debug.SetMemoryLimit(10 * 1024 * 1024)
}

func solution(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Timeout or canceled")
            return
        default:
            // 模拟长时间任务
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    go solution(ctx)

    <-ctx.Done() // 等待超时或完成
    fmt.Println("Main done.")
}
```

## 判题接口调用示例
### 判题沙箱判题接口调用模板
```python
import httpx

url = "http://localhost:8080/sandbox/judgement"
data = {
    "language": "编程语言名称",  # python | java | c | cpp | go
    "question_id": 1,  # 当前题目对应的 ID
    "judge_template": "该编程语言对应的判题模板代码",
    "solution_code": "用户编写的代码",
    "tests": [
        {"test_id": 1, "input_output": "输入输出1"}, 
        {"test_id": 2, "input_output": "输入输出2"}, 
        ...
    ],
    "time_limit": "时间限制",
    "memory_limit": "内存限制",
    "user_id": "用户ID"
}
response = httpx.post(url, json=data)
```

### 接口响应模板
```json
{
  "code": 200,  // 200表示判题成功（不代表通过测试用例），422表示输入数据有问题
  "message": "OK | 其他消息",
  "results": [
    {
      "tets_id": 1, 
      "status": 1,  // 代码运行异常时为-1
      "result": "判题结果（OK | ErrorMessage）",
      "solution": "解题函数的输出（字符串）",
      "criterion": "正确的输出（字符串）",
      "time_consumed": "运行该测试用例消耗的时间",
      "memory_consumed": "运行该测试用例消耗的内存",
      "is_success": true  // 是否通过该测试用例
    },
    ...
  ]
}
```

### 真实调用示例
```python
import httpx

language = "python"
question_id = 1
judge_template = """
from solution_code import solution  # 导入用户编写的解题函数

# 判题模板负责接收输入输出
n = int(input())
nums = map(int, input().split())
criterion = int(input())
# 调用用户编写的解题函数，传入参数并接收返回值
try:
    ans = solution(nums)
except Exception as exception:
    # 处理异常
    pass
else:
    assert ans == criterion
"""
solution_code = """
def solution(nums: lits[int]) -> int:
    return sum(nums)
"""
tests = [
    {
        "test_id": 1,
        "input_output": "3\n1 2 3\n6"
    },
    {
        "test_id": 2,
        "input_output": "2\n10 11\n22"
    }
]
time_limit = 1000 * 20  # 单位是ms
memory_limit = 500  # 单位是MB

data = {
    "language": language,
    "question_id": question_id,
    "judge_template": judge_template,
    "solution_code": solution_code,
    "tests": tests,
    "time_limit": time_limit,
    "memory_limit": memory_limit,
    "user_id": "2025012849102"
}
url = "http://localhost:8080/judgement"
response = httpx.post(url, json=data)
```

### 真实返回响应
```json
{
  "code": 200,
  "message": "OK",
  "results": [
    {
      "tets_id": 1, 
      "status": 1,
      "result": "OK",
      "solution": "6",
      "criterion": "6",
      "time_consumed": 212,
      "memory_consumed": 0.1,
      "is_success": true
    },
    {
      "tets_id": 2,
      "status": 1,
      "result": "OK",
      "solution": "22",
      "criterion": "22",
      "time_consumed": 234,
      "memory_consumed": 0.3,
      "is_success": true
    },
  ]
}
```
