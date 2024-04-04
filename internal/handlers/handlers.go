package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
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

	vm := otto.New()
	result, err:= vm.Run(spec.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result, "input": spec.Code})
}
