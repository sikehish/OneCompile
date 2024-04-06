package utils

import (
	"fmt"
	"os/exec"
)

func CheckImageExists(image string) error {
	cmd := exec.Command("docker", "image", "inspect", "--format='{{.Id}}'", image)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("image %s not found locally", image)
	}
	return nil
}
