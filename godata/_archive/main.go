package main

import (
	"context"
	"fmt"
	"log"

	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

func main() {
	// Create model.
	modelInstance := openai.New("deepseek-chat",
		openai.WithVariant(openai.VariantDeepSeek),
	)

	// Create tool.
	calculatorTool := function.NewFunctionTool(
		calculator,
		function.WithName("calculator"),
		function.WithDescription("Execute addition, subtraction, multiplication, and division. "+
			"Parameters: a, b are numeric values, op takes values add/sub/mul/div; "+
			"returns result as the calculation result."),
	)

	// Enable streaming output.
	genConfig := model.GenerationConfig{
		Stream: true,
	}

	// Create Agent.
	agent := llmagent.New("assistant",
		llmagent.WithModel(modelInstance),
		llmagent.WithTools([]tool.Tool{calculatorTool}),
		llmagent.WithGenerationConfig(genConfig),
	)

	// Create Runner.
	runner := runner.NewRunner("calculator-app", agent)

	// Execute conversation.
	ctx := context.Background()
	events, err := runner.Run(ctx,
		"user-001",
		"session-001",
		model.NewUserMessage("Calculate what 2+3 equals"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Process event stream.
	for event := range events {
		if event.Object == "chat.completion.chunk" {
			fmt.Print(event.Response.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
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
		result = req.A / req.B
	default:
		return calculatorRsp{}, fmt.Errorf("invalid operation: %s", req.Op)
	}
	return calculatorRsp{Result: result}, nil
}

type calculatorReq struct {
	A  float64 `json:"A"  jsonschema:"description=First integer operand,required"`
	B  float64 `json:"B"  jsonschema:"description=Second integer operand,required"`
	Op string  `json:"Op" jsonschema:"description=Operation type,enum=add,enum=sub,enum=mul,enum=div,required"`
}

type calculatorRsp struct {
	Result float64 `json:"result"`
}
