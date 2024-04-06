package utils

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Applies to C/C++

func RunCppInDocker(code string) (string, error) {

	const image = "gcc:latest"
	if err := CheckImageExists(image); err != nil {
		if err := PullDockerImage(image); err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", image, "bash", "-c", fmt.Sprintf("echo '%s' > main.cpp && g++ main.cpp -o main && ./main", code))

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("execution timed out")
		}
		return "", fmt.Errorf("failed to execute C++ code: %v, output: %s", err, output)
	}

	return string(output), nil
}
