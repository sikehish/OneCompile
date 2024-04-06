package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sikehish/OneCompile/internal/utils"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	lang, code := spec.Language, spec.Code

	var output string
	var err error

	switch lang {
	case "js", "javascript":
		output, err = utils.RunJsInDocker(code)
	case "py", "python":
		output, err = utils.RunPythonInDocker(code)
	case "java":
		output, err = utils.RunJavaInDocker(code)
	case "c", "cpp":
		output, err = utils.RunCppInDocker(code)
	case "go", "golang":
		output, err = utils.RunGoInDocker(code)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported language"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": output})
}
