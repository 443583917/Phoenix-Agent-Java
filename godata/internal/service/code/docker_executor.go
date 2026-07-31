package code

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/phoenix-agent-go/internal/config"
)

type DockerExecutor struct {
	config config.CodeExecutorConfig
	logger *zap.Logger
}

func NewDockerExecutor(cfg config.CodeExecutorConfig) *DockerExecutor {
	return &DockerExecutor{config: cfg, logger: zap.L().Named("code.docker")}
}

func (e *DockerExecutor) Execute(ctx context.Context, code string) (*ExecutionResult, error) {
	if e.config.Type == "simulation" {
		return &ExecutionResult{Success: true, Output: "Simulation mode - Docker executor not available"}, nil
	}

	containerName := fmt.Sprintf("%s%d", e.config.ContainerPrefix, time.Now().UnixNano())
	tmpDir, err := os.MkdirTemp("", "phoenix-python-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	scriptPath := filepath.Join(tmpDir, "script.py")
	if err := os.WriteFile(scriptPath, []byte(code), 0o644); err != nil {
		return nil, fmt.Errorf("write script: %w", err)
	}

	timeout := time.Duration(e.config.CodeTimeout) * time.Second
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	args := []string{
		"run", "--rm",
		"--name", containerName,
		"--network", e.config.NetworkMode,
		"--memory", fmt.Sprintf("%dm", e.config.LimitMemoryMB),
		"--cpus", fmt.Sprintf("%d", e.config.CPUCores),
		"-v", fmt.Sprintf("%s:/workspace", tmpDir),
		"-w", "/workspace",
		e.config.ImageName,
		"python", "script.py",
	}

	cmd := exec.CommandContext(execCtx, "docker", args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()

	if execCtx.Err() == context.DeadlineExceeded {
		return &ExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("execution timeout (%ds)", e.config.CodeTimeout),
		}, nil
	}

	result := &ExecutionResult{
		Success: runErr == nil,
		Output:  stdout.String(),
	}
	if runErr != nil {
		result.Error = stderr.String()
		if result.Error == "" {
			result.Error = runErr.Error()
		}
	}

	e.logger.Info("docker execution completed",
		zap.Bool("success", result.Success),
		zap.Int("outputLen", len(result.Output)),
	)

	return result, nil
}
