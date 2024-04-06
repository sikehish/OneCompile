package utils

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func RunJavaInDocker(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "openjdk:latest", "/bin/bash", "-c", fmt.Sprintf("echo '%s' > Main.java && javac Main.java && java Main", code))

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("execution timed out")
		}
		return "", fmt.Errorf("failed to execute Java code: %v, output: %s", err, output)
	}

	return string(output), nil
}
