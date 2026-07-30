package nodes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

type PythonExecuteNode struct{}

func (n *PythonExecuteNode) Name() string {
	return "python_execute"
}

func (n *PythonExecuteNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	code := nl2state.PythonGenerateOutput
	if code == "" || code == "# Python code generation failed" {
		nl2state.PythonExecuteOutput = "无Python代码可执行"
		nl2state.PythonIsSuccess = false
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                       nl2state,
			nl2sqltypes.StateKeyPythonExecuteOutput: "无Python代码可执行",
			nl2sqltypes.StateKeyPythonIsSuccess:     false,
		}, nil
	}

	// Try python3 or python
	pythonBin := "python3"
	if _, err := exec.LookPath("python3"); err != nil {
		pythonBin = "python"
	}

	// Write code to temp file
	tmpFile, err := os.CreateTemp("", "nl2sql-*.py")
	if err != nil {
		nl2state.PythonExecuteOutput = fmt.Sprintf("创建临时文件失败: %s", err.Error())
		nl2state.PythonIsSuccess = false
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                       nl2state,
			nl2sqltypes.StateKeyPythonExecuteOutput: nl2state.PythonExecuteOutput,
			nl2sqltypes.StateKeyPythonIsSuccess:     false,
		}, nil
	}
	defer os.Remove(tmpFile.Name())

	// Prepend imports if needed
	fullCode := code
	if !strings.Contains(code, "import pandas") && !strings.Contains(code, "import numpy") {
		fullCode = "import pandas as pd\nimport numpy as np\nimport json\nimport sys\n\n" + code
	}

	// Add print for JSON output
	if !strings.Contains(code, "print(json.dumps") {
		fullCode += "\n\n# Output results as JSON\nresult = locals().get('result', None)\nif result is not None:\n    print('===RESULT_START===')\n    print(json.dumps(result, ensure_ascii=False, default=str))\n    print('===RESULT_END===')\n"
	}

	if _, err := tmpFile.Write([]byte(fullCode)); err != nil {
		nl2state.PythonExecuteOutput = fmt.Sprintf("写入临时文件失败: %s", err.Error())
		nl2state.PythonIsSuccess = false
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                       nl2state,
			nl2sqltypes.StateKeyPythonExecuteOutput: nl2state.PythonExecuteOutput,
			nl2sqltypes.StateKeyPythonIsSuccess:     false,
		}, nil
	}
	tmpFile.Close()

	// Execute Python
	cmd := exec.CommandContext(ctx, pythonBin, tmpFile.Name())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := fmt.Sprintf("执行失败: %s\nstderr: %s", err.Error(), stderr.String())
		nl2state.PythonExecuteOutput = output
		nl2state.PythonIsSuccess = false
		nl2state.PythonTriesCount = nl2state.PythonTriesCount + 1
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                        nl2state,
			nl2sqltypes.StateKeyPythonExecuteOutput:  output,
			nl2sqltypes.StateKeyPythonIsSuccess:      false,
			nl2sqltypes.StateKeyPythonTriesCount:     nl2state.PythonTriesCount,
		}, nil
	}

	// Capture stderr as additional info
	stderrStr := stderr.String()
	stdoutStr := stdout.String()

	// Extract structured result if present
	_ = stdoutStr
	if start := strings.Index(stdoutStr, "===RESULT_START==="); start >= 0 {
		start += len("===RESULT_START===")
		if end := strings.Index(stdoutStr[start:], "===RESULT_END==="); end >= 0 {
			_ = strings.TrimSpace(stdoutStr[start : start+end])
		}
	}

	output := fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdoutStr, stderrStr)
	if stderrStr != "" {
		output = fmt.Sprintf("执行成功(有stderr输出):\nstdout:\n%s\nstderr:\n%s", stdoutStr, stderrStr)
	} else {
		output = fmt.Sprintf("执行成功:\n%s", stdoutStr)
	}

	nl2state.PythonExecuteOutput = output
	nl2state.PythonIsSuccess = true
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                       nl2state,
		nl2sqltypes.StateKeyPythonExecuteOutput: output,
		nl2sqltypes.StateKeyPythonIsSuccess:     true,
	}, nil
}
