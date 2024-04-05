package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

type Spec struct {
	Language string `json:"language" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Your server's working fine!",
	})
}

func Execute(c *gin.Context) {
	var spec Spec
	if err := c.BindJSON(&spec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := runJsInDocker(spec.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": output})
}

func runJsInDocker(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Set timeout duration (e.g., 10 seconds)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", "node:latest", "node")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start Docker command: %v", err)
	}

	if _, err := stdin.Write([]byte(code)); err != nil {
		return "", fmt.Errorf("failed to write code to stdin: %v", err)
	}

	stdin.Close()

	//wait for command to finish/timeout
	if err := cmd.Wait(); err != nil {
		//check if the error is due to timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("execution timed out")
		}
		return "", fmt.Errorf("failed to execute code: %v, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	return output, nil
}
