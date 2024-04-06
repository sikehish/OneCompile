package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetContainerID(image string) (string, error) {
	var stdout bytes.Buffer
	cmdInspect := exec.Command("docker", "ps", "--filter", "ancestor="+image, "--format={{.ID}}")
	cmdInspect.Stdout = &stdout
	if err := cmdInspect.Run(); err != nil {
		return "", err
	}
	containerID := strings.TrimSpace(stdout.String())
	if containerID == "" {
		return "", fmt.Errorf("container with image %s not found", image)
	}
	return containerID, nil
}
