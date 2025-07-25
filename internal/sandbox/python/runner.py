import json
import sys
import traceback
from typing import Callable, Any, Sequence


run_result = {
    "status": -1,
    "result": "",
    "is_success": False,
    "answer": None,
    "criterion": None,
    "time_consumed": -1,
    "memory_consumed": -1,
}


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


def print_run_result():
    print(json.dumps(run_result))


def run(runner_class: "BaseRunner"):
    try:
        from solution_code import solution
    except:  # 导入解题函数期间产生异常
        run_result['result'] = parse_exception()
        print_run_result()
    else:
        runner_class(solution).run()


class BaseRunner:
    def __init__(self, solution: Callable):
        self.inputs: Sequence[Any] = None
        self.outputs: Sequence[Any] = None
        self.answer: Any | Sequence[Any] = None
        self.criterion: Any | Sequence[Any] = None
        self.solution: Callable = solution

    def process_input_output(self):
        """
        从一个标准的字符串输入中处理输入和输出结果。

        处理完毕以后需要为 self.inputs, self.outputs 这两个属性赋值。
        """
    
    def check_answer(self, answer: Any) -> bool:
        """
        检查解题函数的运行结果。

        检查完毕以后需要调用 set_answer_and_criterion 方法，将解题函数输出与预期输出传到这个方法中。

        Args:
            answer: 用户解题函数的返回值
        
        Returns:
            is_success: 解题函数的返回值是否与预期值相同
        """
    
    def set_answer_and_criterion(self, answer, criterion):
        self.answer = answer
        self.criterion = criterion

    def run(self):
        self.process_input_output()
        status = 1
        result = "OK"
        answer = None
        is_success = False
        try:
            answer = self.solution(*self.inputs)
        except Exception as _:
            result = parse_exception()
            status = -1
        else:
            is_success = self.check_answer(answer)
        run_result["is_success"] = is_success
        run_result["result"] = result
        run_result["status"] = status
        run_result["answer"] = answer
        run_result["criterion"] = self.criterion
        print_run_result()
