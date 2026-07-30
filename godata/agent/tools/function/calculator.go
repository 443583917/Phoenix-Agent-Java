package function

import (
	"context"
	"fmt"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// CalculatorTool returns a function tool that performs arithmetic operations.
// Supports add, sub, mul, div operations on two float64 operands.
func CalculatorTool() tool.Tool {
	return function.NewFunctionTool(
		calculator,
		function.WithName("calculator"),
		function.WithDescription(
			"Execute addition, subtraction, multiplication, and division. "+
				"Parameters: a, b are numeric values, op takes values add/sub/mul/div; "+
				"returns result as the calculation result.",
		),
	)
}

func calculator(ctx context.Context, req calculatorReq) (calculatorRsp, error) {
	var result float64
	switch req.Op {
	case "add", "+":
		result = req.A + req.B
	case "sub", "-":
		result = req.A - req.B
	case "mul", "*":
		result = req.A * req.B
	case "div", "/":
		if req.B == 0 {
			return calculatorRsp{}, fmt.Errorf("division by zero")
		}
		result = req.A / req.B
	default:
		return calculatorRsp{}, fmt.Errorf("invalid operation: %s", req.Op)
	}
	return calculatorRsp{Result: result}, nil
}

type calculatorReq struct {
	A  float64 `json:"A"  jsonschema:"description=First numeric operand,required"`
	B  float64 `json:"B"  jsonschema:"description=Second numeric operand,required"`
	Op string  `json:"Op" jsonschema:"description=Operation type,enum=add,enum=sub,enum=mul,enum=div,required"`
}

type calculatorRsp struct {
	Result float64 `json:"result"`
}
