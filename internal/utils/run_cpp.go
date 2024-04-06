package utils

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Applies to C/C++
func RunCppInDocker(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", "gcc:latest", "bash", "-c", fmt.Sprintf("echo '%s' > main.cpp && g++ main.cpp -o main && ./main", code))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("execution timed out")
		}
		return "", fmt.Errorf("failed to execute C++ code: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
