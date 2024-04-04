package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sikehish/OneCompile/internal/handlers"
)

func newRouter() *gin.Engine {
	r := gin.Default()
	//Testing if the API's working
	r.GET("/test", handlers.TestHandler)
	r.POST("/execute", handlers.Execute)
	return r
}

func RunServer(addr string) error {
	router := newRouter()
	return router.Run(addr)
}
