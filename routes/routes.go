package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup sets all routes & groups used in program.
func Setup(router *gin.Engine) {
	router.Use(cors.Default())
	// router.GET("/apidoc/*any", swagger.WrapHandler(swaggerFiles.Handler))

	userSetup(router)
}
