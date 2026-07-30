package code

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type ExecutorType string

const (
	ExecutorTypeDocker    ExecutorType = "docker"
	ExecutorTypeLocal     ExecutorType = "local"
	ExecutorTypeSimulation ExecutorType = "simulation"
)

type ExecutionResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

type CodePoolExecutor struct {
	executorType ExecutorType
	logger       *zap.Logger
}

func NewCodePoolExecutor(executorType ExecutorType) *CodePoolExecutor {
	return &CodePoolExecutor{
		executorType: executorType,
		logger:       zap.L().Named("code.executor"),
	}
}

func (e *CodePoolExecutor) Execute(ctx context.Context, code string) (*ExecutionResult, error) {
	e.logger.Info("executing code",
		zap.String("type", string(e.executorType)),
		zap.Int("codeLen", len(code)),
	)

	switch e.executorType {
	case ExecutorTypeSimulation:
		return &ExecutionResult{
			Success: true,
			Output:  "Simulation mode - code execution not available",
		}, nil
	case ExecutorTypeLocal:
		return e.executeLocal(ctx, code)
	case ExecutorTypeDocker:
		return e.executeDocker(ctx, code)
	default:
		return nil, fmt.Errorf("unsupported executor type: %s", e.executorType)
	}
}

func (e *CodePoolExecutor) executeLocal(ctx context.Context, code string) (*ExecutionResult, error) {
	e.logger.Warn("local executor not yet implemented, returning simulation result")
	return &ExecutionResult{
		Success: true,
		Output:  "Local executor not yet implemented",
	}, nil
}

func (e *CodePoolExecutor) executeDocker(ctx context.Context, code string) (*ExecutionResult, error) {
	e.logger.Warn("docker executor not yet implemented, returning simulation result")
	return &ExecutionResult{
		Success: true,
		Output:  "Docker executor not yet implemented",
	}, nil
}
