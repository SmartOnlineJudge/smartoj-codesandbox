import json
import traceback
from typing import Callable, Any, Sequence


class RunnerBase:
    def __init__(self, solution: Callable):
        self.inputs: Sequence[Any, ...] = None
        self.outputs: Sequence[Any, ...] = None
        self.answer: Any | Sequence[Any, ...] = None
        self.criterion: Any | Sequence[Any, ...] = None
        self.solution: Callable = solution

    def process_input_output(self):
        pass
    
    def check_answer(self, answer: Any) -> bool:
        pass

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
            result = traceback.format_exc()
            status = -1
        else:
            is_success = self.check_answer(answer)
        run_result = {
            "status": status,
            "result": result,
            "is_success": is_success,
            "answer": self.answer,
            "criterion": self.criterion,
            "time_consumed": 554,
            "memory_consumed": 660,
        }
        print(json.dumps(run_result))
