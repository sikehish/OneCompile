package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunPythonInDocker(code string) (string, error) {

	const image = "python:latest"
	if err := CheckImageExists(image); err != nil {
		if err := PullDockerImage(image); err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", image, "python")

	cmd.Stdin = strings.NewReader(code)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			containerID, err := GetContainerID(image)
			if err != nil {
				return "", fmt.Errorf("execution timed out and failed to retrieve container ID: %v", err)
			}
			if err := exec.Command("docker", "stop", containerID).Run(); err != nil {
				return "", fmt.Errorf("failed to stop container %s: %v", containerID, err)
			}
			return "", fmt.Errorf("execution timed out")
		}
		return "", fmt.Errorf("failed to execute code: %v, output: %s", err, string(output))
	}

	return string(output), nil
}
