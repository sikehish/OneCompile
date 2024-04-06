package utils

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func PullDockerImage(image string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // Set a timeout for pulling the image
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "pull", image)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull Docker image %s: %v", image, err)
	}

	return nil
}
