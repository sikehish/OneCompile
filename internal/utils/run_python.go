package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunPythonInDocker(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", "python:latest", "python")

	cmd.Stdin = strings.NewReader(code)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("execution timed out")
		}
		return "", fmt.Errorf("failed to execute code: %v, output: %s", err, string(output))
	}

	return string(output), nil
}
