package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunJavaInDocker(code string) (string, error) {
	tempDir, err := os.MkdirTemp("", "java")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write the Java source code to a temporary file
	javaFile := filepath.Join(tempDir, "Main.java")
	if err := os.WriteFile(javaFile, []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write Java source file: %v", err)
	}

	if _, err := os.Stat(javaFile); os.IsNotExist(err) {
		return "", fmt.Errorf("java source file does not exist: %s", javaFile)
	}

	cmd := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/code", tempDir), "openjdk:latest", "javac", "/code/Main.java")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to compile Java source code: %v, stderr: %s", err, stderr.String())
	}

	cmd = exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/code", tempDir), "openjdk:latest", "java", "-cp", "/code", "Main")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute Java code: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
