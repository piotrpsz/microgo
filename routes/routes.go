package routes

import (
	"github.com/gin-gonic/gin"
)

// Setup sets all routes & groups used in program.
func Setup(router *gin.Engine) {
	userSetup(router)
}
