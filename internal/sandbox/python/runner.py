import json
import os
import sys
import traceback
import signal
import resource
import time
from typing import Any, Sequence
from contextlib import contextmanager


def parse_exception() -> str:
    exc_type, exc_value, exc_tb = sys.exc_info()
    tb_list = traceback.extract_tb(exc_tb)
    last_frame = tb_list[-1]

    error_msg = f"{exc_type.__name__}: {exc_value}"
    code_line = last_frame.line.strip() if last_frame.line else ""

    return json.dumps({
        "error_msg": error_msg,
        "error_lineno": last_frame.lineno,
        "error_colno": last_frame.colno,
        "error_line": code_line
    })


class ResourceLimiter:
    """
    资源使用限制器

    Args:
        time_limit: 最大执行时间（毫秒）
        memory_limit: 最大内存使用量（MB）
    """
    def __init__(self, time_limit: int, memory_limit: float):
        self.time_limit = int(time_limit / 1000)  # 毫秒 -> 秒
        self.memory_limit = int(memory_limit * 1024 * 1024)
        self.start_time = None
        self.start_memory = None
        self.end_time = None
        self.end_memory = None
        
    @contextmanager
    def limit_resources(self):
        """上下文管理器，用于限制资源并跟踪使用情况"""
        # 记录开始指标
        self.start_time = time.time()
        self.start_memory = self._get_memory_usage()
        
        # 设置时间限制
        signal.signal(signal.SIGALRM, self._timeout_handler)
        signal.alarm(self.time_limit)
        
        # 设置内存限制（虚拟内存限制）
        resource.setrlimit(resource.RLIMIT_AS, (self.memory_limit, self.memory_limit))
        
        try:
            yield
        finally:
            # 记录结束指标
            self.end_time = time.time()
            self.end_memory = self._get_memory_usage()

            # 取消定时器
            signal.alarm(0)
    
    def _timeout_handler(self, signum, frame):
        """处理超时信号"""
        raise TimeoutError
    
    def _get_memory_usage(self):
        """
        获取当前进程内存使用量（字节）
        通过读取 Linux /proc 文件系统获取准确的内存信息
        """
        try:
            # 读取当前进程的内存信息
            with open(f'/proc/{os.getpid()}/status', 'r') as f:
                for line in f:
                    if line.startswith('VmRSS:'):  # 实际使用的物理内存
                        # 格式: VmRSS: 12345 kB
                        return int(line.split()[1])
        except (IOError, ValueError, IndexError):
            # 如果无法读取 /proc 文件，使用resource模块作为备选
            usage = resource.getrusage(resource.RUSAGE_SELF)
            return usage.ru_maxrss
    
    def get_execution_stats(self):
        """
        返回资源使用量

        Returns:
            time_consumed: 时间消耗（毫秒）
            memory_consumed: 内存消耗（MB）
        """
        if self.start_time is None or self.end_time is None:
            return None
            
        time_consumed = int((self.end_time - self.start_time) * 1000)  # 转换为毫秒
        memory_consumed = round((self.end_memory - self.start_memory) / 1024, 2)  # 转换为MB
        
        return time_consumed, memory_consumed


class BaseRunner:
    def __init__(self):
        self.inputs: Sequence[Any] = None
        self.outputs: Sequence[Any] = None
        self.answer: Any | Sequence[Any] = None
        self.time_limit: int = None
        self.memory_limit: float = None

    def process_time_memory_limit(self):
        limits_input = input().split()
        self.time_limit = int(limits_input[0])
        self.memory_limit = float(limits_input[1])

    def process_input_output(self):
        """
        从一个标准的字符串输入中处理输入和输出结果。

        处理完毕以后需要为 self.inputs, self.outputs 这两个属性赋值。
        """

    def check_answer(self, answer: Any) -> bool:
        """
        检查解题函数的运行结果。

        检查完毕以后需要将检查结果和标准结果返回

        Args:
            answer: 用户解题函数的返回值
        
        Returns:
            is_success: 解题函数的返回值是否与预期值相同
            criterion: 正确的结果
        """

    def process_stdin(self):
        """处理标准输入流数据"""
        self.process_time_memory_limit()
        self.process_input_output()

    def run(self):
        self.process_stdin()
        limiter = ResourceLimiter(self.time_limit, self.memory_limit)
        
        status = 1
        result = "OK"
        answer = criterion = None
        is_success = False
        time_consumed, memory_consumed = -1, -1

        with limiter.limit_resources():
            try:
                from solution_code import solution

                answer = solution(*self.inputs)
            except TimeoutError:
                result = "Time Limit Exceeded"
                status = -3
            except MemoryError:
                result = "Memory Limit Exceeded"
                status = -4
            except Exception as _:
                result = parse_exception()
                status = -1

        if status == 1:
            is_success, criterion = self.check_answer(answer)
            time_consumed, memory_consumed = limiter.get_execution_stats()

        run_result = {
            "status": status,
            "result": result,
            "is_success": is_success,
            "answer": str(answer),
            "criterion": str(criterion),
            "time_consumed": time_consumed,
            "memory_consumed": memory_consumed
        }

        print("<SandboxOutput-Start-" + json.dumps(run_result) + "-SandboxOutput-End>")
