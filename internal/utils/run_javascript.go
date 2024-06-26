package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunJsInDocker(code string) (string, error) {

	const image = "node:alpine"
	if err := CheckImageExists(image); err != nil {
		if err := PullDockerImage(image); err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", image, "node")

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

// package utils

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"os/exec"
// 	"time"
// )

// func RunJsInDocker(code string) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Set timeout duration (e.g., 10 seconds)
// 	defer cancel()

// 	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", "node:latest", "node")

// 	var stdout bytes.Buffer
// 	cmd.Stdout = &stdout

// 	var stderr bytes.Buffer
// 	cmd.Stderr = &stderr

// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create stdin pipe: %v", err)
// 	}

// 	if err := cmd.Start(); err != nil {
// 		return "", fmt.Errorf("failed to start Docker command: %v", err)
// 	}

// 	if _, err := stdin.Write([]byte(code)); err != nil {
// 		return "", fmt.Errorf("failed to write code to stdin: %v", err)
// 	}

// 	stdin.Close()

// 	//wait for command to finish/timeout
// 	if err := cmd.Wait(); err != nil {
// 		//check if the error is due to timeout
// 		if ctx.Err() == context.DeadlineExceeded {
// 			return "", fmt.Errorf("execution timed out")
// 		}
// 		return "", fmt.Errorf("failed to execute code: %v, stderr: %s", err, stderr.String())
// 	}

// 	output := stdout.String()
// 	return output, nil
// }
