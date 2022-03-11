package routers

import (
	controllers "cve-sa-backend/controllers/manage"
	"github.com/gin-gonic/gin"
)

func ManageRouters(r *gin.Engine) {

	manager := r.Group("/cve-security-notice-server")
	{
		manager.Use(managerAuth)
		manager.POST("/upload", controllers.Upload)
		manager.POST("/deleteCVE", controllers.DeleteCVE)
	}
}
