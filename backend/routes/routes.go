package routes

import (
	"platform/backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize controllers consistently
	vmController := &controllers.VMController{}
	// Minimal public vm routes (no auth)
	r.POST("/vm/start", vmController.Startvm)
	r.POST("/vm/:id/stop", vmController.Stopvm)
	r.GET("/vm/:id/stream", vmController.Streamvm)
	r.GET("/vm/:id/info", vmController.Info)
}
