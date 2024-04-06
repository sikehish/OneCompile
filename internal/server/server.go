package server

import (
	"github.com/gin-contrib/cors" // Import the cors package
	"github.com/gin-gonic/gin"
	"github.com/sikehish/OneCompile/internal/handlers"
)

func newRouter() *gin.Engine {
	r := gin.Default()

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	// Testing if the API's working
	r.GET("/test", handlers.TestHandler)
	r.POST("/execute", handlers.Execute)
	return r
}

func RunServer(addr string) error {
	router := newRouter()

	return router.Run(addr)
}
