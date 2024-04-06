package utils

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

func RunGoInDocker(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "golang:latest", "/bin/bash", "-c", fmt.Sprintf("echo '%s' > main.go && go run main.go", code))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("execution timed out")
		}
		return "", fmt.Errorf("failed to execute Go code: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
