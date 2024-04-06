package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func RunCppInDocker(code string) (string, error) {
	tempDir, err := os.MkdirTemp("", "cpp")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cppFile := filepath.Join(tempDir, "main.cpp")
	if err := os.WriteFile(cppFile, []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write C++ source file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-v", fmt.Sprintf("%s:/code", tempDir), "gcc:latest", "/bin/bash", "-c", "g++ /code/main.cpp -o /code/main && /code/main")
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
