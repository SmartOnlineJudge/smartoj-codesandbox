package sandbox

import (
	"fmt"

	"smartoj-codesandbox/internal/types"
)


func executePython(jd *types.JudgementData, results *types.Results, workspace string) int {
	for _, test := range jd.Tests {
		fmt.Println(test.InputOutput)
	}
	return -1
}
